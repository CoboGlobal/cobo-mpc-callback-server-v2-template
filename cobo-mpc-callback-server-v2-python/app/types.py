from dataclasses import dataclass
from typing import Optional


# Constants
class Action:
    APPROVE = "APPROVE"
    REJECT = "REJECT"


class Status:
    OK = 0
    INVALID_REQUEST = 10
    INVALID_TOKEN = 20
    INTERNAL_ERROR = 30


@dataclass
class PackageDataClaim:
    """JWT claims with package data and standard JWT claims"""

    package_data: Optional[str] = None
    aud: Optional[str] = None  # Audience
    exp: Optional[int] = None  # Expiration Time
    jti: Optional[str] = None  # JWT ID
    iat: Optional[int] = None  # Issued At
    iss: Optional[str] = None  # Issuer
    nbf: Optional[int] = None  # Not Before
    sub: Optional[str] = None  # Subject

    def to_dict(self):
        claims = {
            "package_data": self.package_data,
            "aud": self.aud,
            "exp": self.exp,
            "jti": self.jti,
            "iat": self.iat,
            "iss": self.iss,
            "nbf": self.nbf,
            "sub": self.sub,
        }
        return {k: v for k, v in claims.items() if v is not None}


@dataclass
class Request:
    """Request structure"""

    request_id: Optional[str] = None
    request_type: Optional[int] = None
    request_detail: Optional[str] = None
    extra_info: Optional[str] = None

    def to_dict(self):
        return {
            "request_id": self.request_id,
            "request_type": self.request_type,
            "request_detail": self.request_detail,
            "extra_info": self.extra_info,
        }


@dataclass
class Response:
    """Response structure"""

    status: int = Status.OK
    request_id: Optional[str] = None
    action: Optional[str] = None  # [APPROVE, REJECT]
    err_str: Optional[str] = None

    def __str__(self):
        return f"Status: {self.status}, RequestID: {self.request_id}, Action: {self.action}, ErrStr: {self.err_str}"

    def to_dict(self):
        return {
            "status": self.status,
            "request_id": self.request_id,
            "action": self.action,
            "error": self.err_str,
        }
