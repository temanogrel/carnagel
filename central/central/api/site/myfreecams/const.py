from enum import Enum


class VideoStates(Enum):
    TX_IDLE = 0
    TX_RESET = 1
    TX_AWAY = 2
    TX_CONFIRMING = 11
    TX_PVT = 12
    TX_GRP = 13
    TX_KILL_MODEL = 15
    RX_IDLE = 90
    RX_PVT = 91
    RX_VOY = 92
    RX_GRP = 93
    OFFLINE = 127


class AccessLevel(Enum):
    GUEST = 0
    BASIC = 1
    PREMIUM = 2
    MODEL = 4
    ADMIN = 5