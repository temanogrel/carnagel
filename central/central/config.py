import common.rabbitmq

# Celery
BROKER_URL = common.rabbitmq.rabbitmq_dsn('central')

CELERY_TASK_SERIALIZER = 'json'
CELERY_ACCEPT_CONTENT = ['json']
CELERY_RESULT_SERIALIZER = 'json'
CELERY_DISABLE_RATE_LIMITS = True
CELERY_ENABLE_UTC = True
CELERY_ROUTES = {
    'central.tasks.publish_recording': {
        'queue': 'publish'
    },

    'central.tasks.process_death_file': {
        'queue': 'death_file'
    },

    'download.tasks.download_myfreecams': {
        'queue': 'downloading',
    },

    'download.tasks.download_chaturbate': {
        'queue': 'downloading',
    },

    'transcode.tasks.transcode': {
        'queue': 'transcode'
    },

    'storage.tasks.upload_media': {
        'exchange': 'storage_uploads',
        'exchange_type': 'direct'
    },

    'storage.tasks.publish_video': {
        'exchange': 'storage_uploads',
        'exchange_type': 'direct'
    },

    'storage.tasks.reupload_image': {
        'exchange': 'storage_uploads',
        'exchange_type': 'direct'
    },

    'storage.tasks.reupload_video': {
        'exchange': 'storage_uploads',
        'exchange_type': 'direct'
    }
}

# MyFreeCams configuration

MFC_PING_INTERVAL = 30
MFC_INTERSECT_INTERVAL = 10
MFC_RELOAD_MODELS_INTERVAL = 180

# Chaturbate configuration

CBC_SECTIONS = (
    'female',
    'male',
    'couple'
)

CBC_CREDENTIALS = {
    'yilltoestryu8': '$5nf9V!uW0Fk@1iEma',
    'juniper589ol': '$5nf9V!uW0Fk@1iEma',
    'tragicalien75': '$5nf9V!uW0Fk@1iEma',
    'monstertrucker46': '$5nf9V!uW0Fk@1iEma',
    'yungstud005': '$5nf9V!uW0Fk@1iEma',
    'rasputin61': '$5nf9V!uW0Fk@1iEma',
	'princeton709432': '$5nf9V!uW0Fk@1iEma',
	'appleby781774': '$5nf9V!uW0Fk@1iEma',
	'hunglikea653992': '$5nf9V!uW0Fk@1iEma',
	'jimmythebean9034466': '$5nf9V!uW0Fk@1iEma',
	'wonderboy704119': '$5nf9V!uW0Fk@1iEma',
	'racerforyou8995521': '$5nf9V!uW0Fk@1iEma',
	'kingtroy7884012': '$5nf9V!uW0Fk@1iEma',
	'askingyoukiss78499922': '$5nf9V!uW0Fk@1iEma',
	'uptimesinger0019953': '$5nf9V!uW0Fk@1iEma'
}

CBC_SCAN_INTERVAL = 60
