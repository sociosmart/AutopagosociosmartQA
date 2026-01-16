import sentry_sdk
import uuid
from flask import Blueprint, current_app, g
from flask_apispec import doc, marshal_with, use_kwargs

from internal.db import db
from internal.dtos.common import AuthorizationAppHeaders, GeneralMessage
from internal.dtos.payments import (
    MakeAPaymentRequest,
    MakeAPaymentResponse,
    PaymentConfirmationPathParams,
    PaymentConfirmationRequest,
)
from internal.enums.deposits import MovementType
from internal.enums.payments import PaymentStatus
from internal.exceptions.http import GenericError, InternalServerError, NotFoundError
from internal.middlewares.authorized_app import authorized_app_auth
from internal.models.deposits import Movement
from internal.models.payments import Payment, PaymentFundsTrack
from internal.repositories.deposits import (
    bulk_deposit_update_by_id,
    create_movement,
    get_balance_by_customer,
    get_deposits_with_funds,
)
from internal.repositories.customers import get_customer_by_external_id
from internal.repositories.gas_stations import get_legal_name_by_external_id
from internal.repositories.gift_cards import (
    get_gift_card_by_card_key,
    update_gift_card_by_id,
)
from internal.repositories.payments import (
    bulk_create_funds_track,
    check_if_gift_card_is_already_used,
    create_payment,
    get_payment_by_id,
    update_payment_by_id,
)


bp = Blueprint("payments", __name__)


@bp.post("/")
@doc(tags=["Payments"], description="Make a Payment")
@use_kwargs(AuthorizationAppHeaders, location="headers")
@authorized_app_auth()
@use_kwargs(MakeAPaymentRequest, location="json")
@marshal_with(MakeAPaymentResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(
    GeneralMessage,
    code="406",
    description="Unsufficient funds for customer or given gift card",
)
@marshal_with(
    GeneralMessage,
    code="404",
    description="Legal name, gift card or customer does not exist",
)
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
@marshal_with(GeneralMessage, code="409", description="Gift Card in use")
def make_a_payment(
    external_legal_name_id,
    amount,
    external_customer_id=None,
    card_key=None,
    **kwargs,
):
    app = g.authorized_app
    customer = None
    gift_card = None
    try:
        legal_name = get_legal_name_by_external_id(external_legal_name_id)
    except Exception as e:
        current_app.logger.error(f"Error while trying to bring legal name data - {e}")
        raise InternalServerError
    else:
        if not legal_name:
            raise NotFoundError("Legal name not found")

    if external_customer_id:
        try:
            customer = get_customer_by_external_id(external_customer_id)
        except Exception as e:
            current_app.logger.error(f"Error while trying to bring customer data - {e}")
            raise InternalServerError
        else:
            if not customer:
                raise NotFoundError("Given external customer id not found")

        try:
            balance = get_balance_by_customer(customer.id, legal_name_id=legal_name.id)
        except Exception as e:
            current_app.logger.error(
                f"Error while trying to bring customer balance - {e}"
            )
            raise InternalServerError
        else:
            if amount > balance["total"]:
                raise GenericError("Not enough funds", 406)
    else:
        try:
            gift_card = get_gift_card_by_card_key(card_key)
        except Exception as e:
            current_app.logger.error(
                f"Error while trying to bring gift card data - {e}"
            )
            raise InternalServerError
        else:
            if not gift_card:
                raise NotFoundError(
                    "Unable to find the gift card with provided card key"
                )
            if amount > gift_card.amount:
                raise GenericError(
                    "This card does not have that credit to be redeemed", 406
                )

        # Checking gift card is not previously used
        try:
            used_card = check_if_gift_card_is_already_used(gift_card.id)
        except Exception as e:
            current_app.logger.error(
                f"Error while trying to bring gift card data - {e}"
            )
            raise InternalServerError
        else:
            if used_card:
                raise GenericError(f"This gift card is already being in use", 409)

    # Registering payment in db

    try:
        payment = create_payment(
            Payment(
                customer=customer,
                gift_card=gift_card,
                legal_name_id=legal_name.id,
                amount=amount,
                created_by=app,
            ),
            auto_commit=False,
        )
        db.session.flush()
        customer_id = None
        if customer:
            customer_id = customer.id
        else:
            customer_id = gift_card.customer_id
        create_movement(
            Movement(
                type=MovementType.FUNDS_RESERVED,
                payment_id=payment.id,
                amount=amount,
                description=f"Fondos reservado por {amount}",
                customer_id=customer_id,
            )
        )
    except Exception as e:
        db.session.rollback()
        current_app.logger.error(f"Error while trying to Register payment in db - {e}")
        raise InternalServerError

    return MakeAPaymentResponse().dump(payment)


@bp.put("/cancelation/<uuid:id>")
@doc(
    tags=["Payments"],
    description="Payment Cancelation",
    params={
        "id": {
            "description": "Payment ID",
            "example": "c35653af-5171-4450-ad1b-d2e39f91c720",
        }
    },
)
@use_kwargs(AuthorizationAppHeaders, location="headers")
@authorized_app_auth()
@marshal_with(GeneralMessage, code="200", description="Payment canceled")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
@marshal_with(GeneralMessage, code="404", description="Payment not found")
def payment_cancelation(id, **kwargs):
    app = g.authorized_app
    # id = uuid.UUID(id)
    try:
        payment = get_payment_by_id(id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while trying {e}")
        raise InternalServerError

    if not payment:
        raise NotFoundError("Payment not found")

    if payment.status != PaymentStatus.FUNDS_RESERVED:
        raise GenericError("Payment already confirmed or cancelled", 412)

    try:
        customer_id = None
        if payment.paid_by_customer:
            customer_id = payment.customer_id
        else:
            customer_id = payment.gift_card.customer_id
        update_payment_by_id(
            payment.id,
            values={"canceled_by_id": app.id, "status": PaymentStatus.CANCELLED},
            auto_commit=False,
        )
        create_movement(
            Movement(
                type=MovementType.PAYMENT_CANCELED,
                amount=payment.amount,
                customer_id=customer_id,
                payment_id=payment.id,
                description="Pago cancelado",
            ),
            auto_commit=False,
        )
        db.session.commit()
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong while trying to cancel the payment - {e}"
        )
        raise InternalServerError

    return GeneralMessage().dump({"message": "canceled"})


@bp.post("/confirmation/<uuid:id>")
@doc(
    tags=["Payments"],
    description="Payment Confirmation",
    params={
        "id": {
            "description": "Payment ID",
            "example": "c35653af-5171-4450-ad1b-d2e39f91c720",
        }
    },
)
@use_kwargs(AuthorizationAppHeaders, location="headers")
@authorized_app_auth()
@use_kwargs(PaymentConfirmationRequest, location="json")
@marshal_with(GeneralMessage, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
@marshal_with(GeneralMessage, code="404", description="Payment not found")
@marshal_with(
    GeneralMessage, code="412", description="Payment already confirmed or canceled"
)
@marshal_with(
    GeneralMessage,
    code="406",
    description="Given confirmation amount greater than reserved",
)
def payment_confirmation(id, amount, **kwargs):
    app = g.authorized_app
    # id = uuid.UUID(id)
    try:
        payment = get_payment_by_id(id)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong while trying to get payment - {e}"
        )
        raise InternalServerError

    if not payment:
        raise NotFoundError("Payment not found")

    if payment.status != PaymentStatus.FUNDS_RESERVED:
        raise GenericError("Payment already confirmed or cancelled", 412)

    if amount > payment.amount:
        raise GenericError("Given amount greater than funds reserved", 406)

    # Deposits where funds will be taken
    try:
        customer_id = payment.customer_id
        if payment.gift_card:
            customer_id = payment.gift_card.customer_id

        deposits = get_deposits_with_funds(customer_id, payment.legal_name_id)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong while trying to get deposits - {e}"
        )
        raise InternalServerError

    to_discount = amount
    deposits_to_update = []
    deposits_track = []
    for deposit in deposits:
        if deposit.difference >= to_discount:
            deposits_to_update.append(
                {"id": deposit.id, "amount_used": deposit.amount_used + to_discount}
            )
            deposits_track.append(
                dict(
                    payment_id=payment.id,
                    deposit_id=deposit.id,
                    amount=to_discount,
                    prev_amount=deposit.difference,
                )
            )
            to_discount = 0
            break

        to_discount = to_discount - deposit.difference
        deposits_to_update.append({"id": deposit.id, "amount_used": deposit.amount})
        deposits_track.append(
            dict(
                payment_id=payment.id,
                deposit_id=deposit.id,
                amount=deposit.difference,
                prev_amount=deposit.difference,
            )
        )
        if to_discount == 0:
            break

    if to_discount > 0:
        sentry_sdk.set_tag("payment_id", str(payment.id))
        sentry_sdk.set_extra("authorized_app", app.id)
        sentry_sdk.set_extra("amount", amount)
        sentry_sdk.capture_message(
            f"Payment could not be totally paid since something happend with deposits, manually intervention."
        )

    # all tracking
    try:
        bulk_deposit_update_by_id(deposits_to_update, auto_commit=False)
        update_payment_by_id(
            payment.id,
            values={
                "status": PaymentStatus.CONFIRMED,
                "amount_confirmed": amount - to_discount,
                "updated_by_id": app.id,
            },
            auto_commit=False,
        )
        create_movement(
            Movement(
                type=MovementType.FUNDS_CONFIRMATION,
                payment_id=payment.id,
                amount=amount,
                description=f"Fondos confirmados por {amount - to_discount}",
                customer_id=customer_id,
            ),
            auto_commit=False,
        )

        if not payment.paid_by_customer:
            gift_card_id = payment.gift_card_id
            update_gift_card_by_id(
                gift_card_id,
                values={"amount_used": amount - to_discount, "redeemed": True},
                auto_commit=False,
            )
            create_movement(
                Movement(
                    type=MovementType.GIFT_CARD_REDEMPTION,
                    gift_card_id=gift_card_id,
                    amount=amount - to_discount,
                    customer_id=customer_id,
                    description=f"Gift card redeemed for {amount - to_discount}",
                ),
                auto_commit=False,
            )

        bulk_create_funds_track(deposits_track, auto_commit=False)
        db.session.commit()
    except Exception as e:
        db.session.rollback()
        current_app.logger.error(
            f"Unable to save payment tracking, manaual action. - {e}"
        )
        raise InternalServerError

    return GeneralMessage().dump({"message": "Confirmed"})
