from ..celery import app


@app.task(name='transcode.tasks.transcode')
def transcode(recording_id: str, auto_upload=True):
    pass
