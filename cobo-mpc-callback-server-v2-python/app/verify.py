import json
import logging
from abc import ABC, abstractmethod
from typing import Optional

from app.request import (
    KeyGenDetail,
    KeyGenRequestInfo,
    KeyReshareDetail,
    KeyReshareRequestInfo,
    KeySignDetail,
    KeySignRequestInfo,
    RequestType,
)
from app.types import Request

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)


class Verifier(ABC):
    """Abstract base class for verifiers"""

    @abstractmethod
    def verify(self, request: Request) -> Optional[str]:
        pass


class TssVerifier(Verifier):
    """TSS request verifier implementation"""

    @classmethod
    def new(cls) -> "TssVerifier":
        return cls()

    @classmethod
    def verify(cls, request: Request) -> Optional[str]:
        """Verify TSS request"""
        if not request:
            return "request is nil"

        try:
            request_type = RequestType(request.request_type)

            if request_type == RequestType.TYPE_PING:
                logger.debug("Got ping request")
                return None
            elif request_type == RequestType.TYPE_KEY_GEN:
                return cls.handle_key_gen(request.request_detail, request.extra_info)
            elif request_type == RequestType.TYPE_KEY_SIGN:
                return cls.handle_key_sign(request.request_detail, request.extra_info)
            elif request_type == RequestType.TYPE_KEY_RESHARE:
                return cls.handle_key_reshare(
                    request.request_detail, request.extra_info
                )
            else:
                return f"not support to process request type {request.request_type}"
        except ValueError as e:
            return f"invalid request type: {e}"
        except Exception as e:
            logger.error(f"Verify error: {str(e)}")
            return str(e)

    @classmethod
    def handle_key_gen(cls, request_detail: str, extra_info: str) -> Optional[str]:
        """Handle key generation request"""
        if not request_detail or not extra_info:
            return "request detail or extra info is empty"

        try:
            logger.debug(
                f"key gen original detail:\n request detail: {request_detail}\nrequest info:\n{extra_info}"
            )

            key_gen_detail = KeyGenDetail.from_json(request_detail)
            request_info = KeyGenRequestInfo.from_json(extra_info)

            logger.debug(
                f"key gen class detail:\n request detail: {key_gen_detail}\nrequest info:\n{request_info}"
            )

            # key gen logic add here

            return None

        except json.JSONDecodeError as e:
            return f"failed to parse key gen json: {str(e)}"
        except Exception as e:
            return f"failed to handle key gen: {str(e)}"

    @classmethod
    def handle_key_sign(cls, request_detail: str, extra_info: str) -> Optional[str]:
        """Handle key signing request"""
        if not request_detail or not extra_info:
            return "request detail or extra info is empty"

        try:
            logger.debug(
                f"key sign original detail:\n request detail: {request_detail}\nrequest info:\n{extra_info}"
            )

            key_sign_detail = KeySignDetail.from_json(request_detail)
            request_info = KeySignRequestInfo.from_json(extra_info)

            logger.debug(
                f"key sign class detail:\n{key_sign_detail}\nrequest info:\n{request_info}"
            )

            # key sign logic add here

            return None

        except json.JSONDecodeError as e:
            return f"failed to parse key sign json: {str(e)}"
        except Exception as e:
            return f"failed to handle key sign: {str(e)}"

    @classmethod
    def handle_key_reshare(cls, request_detail: str, extra_info: str) -> Optional[str]:
        """Handle key reshare request"""
        if not request_detail or not extra_info:
            return "request detail or extra info is empty"

        try:
            logger.debug(
                f"key reshare original detail:\n request detail: {request_detail}\nrequest info:\n{extra_info}"
            )

            key_reshare_detail = KeyReshareDetail.from_json(request_detail)
            request_info = KeyReshareRequestInfo.from_json(extra_info)

            logger.debug(
                f"key reshare class detail:\n{key_reshare_detail}\nrequest info:\n{request_info}"
            )

            # key reshare logic add here

            return None

        except json.JSONDecodeError as e:
            return f"failed to parse key reshare json: {str(e)}"
        except Exception as e:
            return f"failed to handle key reshare: {str(e)}"
