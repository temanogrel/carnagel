import enum

TERMINATION_STATUSES = (
    'NetStream.Play.UnpublishNotify',
    'NetStream.Play.StreamNotFound',
    'NetStream.Play.Failed'
)


class RecordingState(enum.Enum):
    ERROR = 'error'
    SUCCESS = 'success'
    INITIATED = 'initiated'
    RECORDING = 'recording'
    NOT_AVAILABLE = 'not_available'


class WeAreBannedError(Exception):
    pass


class AuthenticationFailureError(Exception):
    pass


class RecordingTerminatedError(Exception):
    def __init__(self, reason: str):
        super().__init__()

        self.reason = reason


class RecordingTimeoutError(RecordingTerminatedError):
    pass


class ModelNotAvailableError(RecordingTerminatedError):
    pass
