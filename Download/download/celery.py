from celery import Celery

app = Celery('download')
app.config_from_object('download.config')
