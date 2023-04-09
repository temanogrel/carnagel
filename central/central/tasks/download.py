from ..celery import app


@app.task(name='download.tasks.download_myfreecams')
def download_mfc_stream(performer_id: int, session_id: int, auto_transcode=True):
    """
    Proxy function for the real process

    :return:
    """
    pass


@app.task(name='download.tasks.download_chaturbate')
def download_cbc_stream(performer_id: int, username: str, api_token: str, auto_transcode=True):
    """
    Proxy function for the real process

    :return:
    """
    pass
