from sqlalchemy import func

from internal.models.deposits import Deposit
from internal.models.payments import Payment

subtotal = func.coalesce(
    func.sum(Deposit.amount - Deposit.amount_used).label("total"), 0.0
)
total = subtotal
reserved_total = func.coalesce(func.sum(Payment.amount).label("total_reserved"), 0.0)
