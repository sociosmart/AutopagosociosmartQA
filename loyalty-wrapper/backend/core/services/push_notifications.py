import httpx

from core.config import settings
from core.logging.logger import logger


class PushNotificationService:

    def send_notification(self, title, message, tokens=[]):
        uri = f"{settings.push.push_notifications_url}/api/notification"
        try:
            response = httpx.post(
                uri,
                headers={
                    "X-APP-KEY": settings.push.push_app_key,
                    "X-API-KEY": settings.push.push_api_key,
                },
                json={
                    "notifications": [
                        {
                            "title": title,
                            "message": message,
                            "platform": 2,
                            "tokens": tokens,
                        }
                    ]
                },
            )

            if response.status_code != 200:
                logger.error(
                    f"Something went wrong while sending push notification - {response.status_code}"
                )
                raise Exception("Unable to send push notification")
        except Exception as e:
            logger.error(f"Something went wrong - {e}")
            raise e


push_notifications_service = PushNotificationService()
