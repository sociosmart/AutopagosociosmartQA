import datetime

from sqlalchemy import DateTime, String
from sqlalchemy.orm import Mapped, mapped_column

from api.internal.db import db


class ScheduledCampaigns(db.Model):
    __tablename__ = "scheduled_campaigns"

    id: Mapped[int] = mapped_column(primary_key=True)
    external_campaign_id: Mapped[str] = mapped_column(String(50))
    scheduled_at: Mapped[datetime.datetime] = mapped_column(DateTime)
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
