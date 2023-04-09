import os
import common.rabbitmq

####################################
### Generic worker configuration ###
####################################

DOWNLOAD_PATH = os.environ.get('DOWNLOAD_PATH', '/opt/downloads/')
HOSTNAME = os.environ.get('HOSTNAME')

############################
### CELERY CONFIGURATION ###
############################

broker_url = common.rabbitmq.rabbitmq_dsn('downloader')
broker_pool_limit = 50

celery_acks_late = True

task_ignore_result = True

worker_prefetch_multiplier = 1
worker_max_tasks_per_child = 50

task_routes = {
    'download.tasks.download_myfreecams': {
        'queue': 'downloading',
    },

    'download.tasks.download_chaturbate': {
        'queue': 'downloading',
    },

    'transcode.tasks.transcode': {
        'queue': 'transcode',
    }
}
