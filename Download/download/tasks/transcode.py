from download.celery import app


@app.task(name='transcode.tasks.transcode')
def transcode(recording_id: int, auto_upload=True):
    """
    Proxy method for the real transcoding task

    :param recording_id:
    :param encoding:

    :return:
    """
    pass
