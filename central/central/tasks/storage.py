from central.celery import app

@app.task(name='storage.tasks.publish_video')
def publish_video(recording_id: str, storage_path: str):
    pass


@app.task(name='storage.tasks.update_metadata')
def update_metadata(recording_id: str, storage_path: str):
    pass


@app.task(name='storage.tasks.delete_file')
def delete_file(path: str):
    pass


@app.task(name='storage.tasks.reupload_image')
def reupload_recording_image(recording_id: str):
    pass


@app.task(name='storage.tasks.reupload_video')
def reupload_recording_video(recording_id: str):
    pass


@app.task(name='storage.tasks.upload_media')
def upload_media(recording_id: int, auto_publish=True):
    pass
