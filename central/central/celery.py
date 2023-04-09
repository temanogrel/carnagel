from celery import Celery

app = Celery('pacific_artifacts')
app.config_from_object('central.config')
