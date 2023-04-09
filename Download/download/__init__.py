# This will make sure the app is always imported when
# Django starts so that shared_task will use this app.
from .celery import app as celery_app
from .tasks.download import download_chaturbate, download_myfreecams

__all__ = ['celery_app']
