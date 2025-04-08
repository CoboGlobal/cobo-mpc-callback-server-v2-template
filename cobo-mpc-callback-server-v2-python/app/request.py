import json
from dataclasses import dataclass
from enum import IntEnum
from typing import List, Optional

import cobo_waas2


# =================== request_type ===================
class RequestType(IntEnum):
    TYPE_PING = 0
    TYPE_KEY_GEN = 1
    TYPE_KEY_SIGN = 2
    TYPE_KEY_RESHARE = 3

    def __str__(self):
        if self == RequestType.TYPE_PING:
            return "Ping"
        elif self == RequestType.TYPE_KEY_GEN:
            return "KeyGen"
        elif self == RequestType.TYPE_KEY_SIGN:
            return "KeySign"
        elif self == RequestType.TYPE_KEY_RESHARE:
            return "KeyReshare"
        return "Unknown"


# =================== request_detail ===================
class CurveType(IntEnum):
    SECP256K1 = 0
    ED25519 = 2


class SignatureType(IntEnum):
    UNKNOWN_TYPE = 0
    ECDSA = 1
    EDDSA = 2
    SCHNORR = 3


class TssProtocol(IntEnum):
    UNKNOWN_PROTOCOL = 0
    GG18 = 1
    LINDELL = 2
    EDDSA_TSS = 3


@dataclass
class KeyGenDetail:
    """Key generation detail"""

    threshold: int
    curve: CurveType
    node_ids: List[str]
    task_id: str
    biz_task_id: str

    def __str__(self):
        return json.dumps(
            {
                "threshold": self.threshold,
                "curve": self.curve,
                "node_ids": self.node_ids,
                "task_id": self.task_id,
                "biz_task_id": self.biz_task_id,
            }
        )

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeyGenDetail instance from a JSON string.

        Args:
            json_str: JSON string containing the key generation detail

        Returns:
            KeyGenDetail: An instance with parsed data
        """
        if not json_str:
            return cls(
                threshold=0,
                curve=CurveType.SECP256K1,
                node_ids=[],
                task_id="",
                biz_task_id="",
            )

        try:
            data = json.loads(json_str)
            return cls(
                threshold=data.get("threshold"),
                curve=CurveType(data.get("curve")),
                node_ids=data.get("node_ids", []),
                task_id=data.get("task_id", ""),
                biz_task_id=data.get("biz_task_id", ""),
            )
        except json.JSONDecodeError:
            return cls(
                threshold=0,
                curve=CurveType.SECP256K1,
                node_ids=[],
                task_id="",
                biz_task_id="",
            )

    """Key generation detail"""


@dataclass
class KeySignDetail:
    """Key signing detail"""

    group_id: str
    root_pub_key: str
    used_node_ids: List[str]
    bip32_path_list: List[str]
    msg_hash_list: List[str]
    tweak_list: List[str]
    signature_type: SignatureType
    tss_protocol: TssProtocol
    task_id: str
    biz_task_id: str

    def __str__(self):
        return json.dumps(
            {
                "group_id": self.group_id,
                "root_pub_key": self.root_pub_key,
                "used_node_ids": self.used_node_ids,
                "bip32_path_list": self.bip32_path_list,
                "msg_hash_list": self.msg_hash_list,
                "tweak_list": self.tweak_list,
                "signature_type": self.signature_type,
                "tss_protocol": self.tss_protocol,
                "task_id": self.task_id,
                "biz_task_id": self.biz_task_id,
            }
        )

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeySignDetail instance from a JSON string.

        Args:
            json_str: JSON string containing the key signing detail

        Returns:
            KeySignDetail: An instance with parsed data
        """
        if not json_str:
            return cls(
                group_id="",
                root_pub_key="",
                used_node_ids=[],
                bip32_path_list=[],
                msg_hash_list=[],
                tweak_list=[],
                signature_type=SignatureType.UNKNOWN_TYPE,  # Default value
                tss_protocol=TssProtocol.UNKNOWN_PROTOCOL,  # Default value
                task_id="",
                biz_task_id="",
            )

        try:
            data = json.loads(json_str)
            return cls(
                group_id=data.get("group_id", ""),
                root_pub_key=data.get("root_pub_key", ""),
                used_node_ids=data.get("used_node_ids", []),
                bip32_path_list=data.get("bip32_path_list", []),
                msg_hash_list=data.get("msg_hash_list", []),
                tweak_list=data.get("tweak_list", []),
                signature_type=SignatureType(
                    data.get("signature_type", SignatureType.UNKNOWN_TYPE.value)
                ),
                tss_protocol=TssProtocol(
                    data.get("tss_protocol", TssProtocol.UNKNOWN_PROTOCOL.value)
                ),
                task_id=data.get("task_id", ""),
                biz_task_id=data.get("biz_task_id", ""),
            )
        except json.JSONDecodeError:
            return cls(
                group_id="",
                root_pub_key="",
                used_node_ids=[],
                bip32_path_list=[],
                msg_hash_list=[],
                tweak_list=[],
                signature_type=SignatureType.UNKNOWN_TYPE,
                tss_protocol=TssProtocol.UNKNOWN_PROTOCOL,
                task_id="",
                biz_task_id="",
            )


@dataclass
class KeyReshareDetail:
    """Key reshare detail"""

    old_group_id: str
    root_pub_key: str
    curve: CurveType
    used_node_ids: List[str]
    old_threshold: int
    new_threshold: int
    new_node_ids: List[str]
    task_id: str
    biz_task_id: str

    def __str__(self):
        return json.dumps(
            {
                "old_group_id": self.old_group_id,
                "root_pub_key": self.root_pub_key,
                "curve": self.curve,
                "used_node_ids": self.used_node_ids,
                "old_threshold": self.old_threshold,
                "new_threshold": self.new_threshold,
                "new_node_ids": self.new_node_ids,
                "task_id": self.task_id,
                "biz_task_id": self.biz_task_id,
            }
        )

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeyReshareDetail instance from a JSON string.

        Args:
            json_str: JSON string containing the key reshare detail

        Returns:
            KeyReshareDetail: An instance with parsed data
        """
        if not json_str:
            return cls(
                old_group_id="",
                root_pub_key="",
                curve=CurveType.SECP256K1,  # Default curve type
                used_node_ids=[],
                old_threshold=0,
                new_threshold=0,
                new_node_ids=[],
                task_id="",
                biz_task_id="",
            )

        try:
            data = json.loads(json_str)
            return cls(
                old_group_id=data.get("old_group_id", ""),
                root_pub_key=data.get("root_pub_key", ""),
                curve=CurveType(data.get("curve", CurveType.SECP256K1.value)),
                used_node_ids=data.get("used_node_ids", []),
                old_threshold=data.get("old_threshold", 0),
                new_threshold=data.get("new_threshold", 0),
                new_node_ids=data.get("new_node_ids", []),
                task_id=data.get("task_id", ""),
                biz_task_id=data.get("biz_task_id", ""),
            )
        except json.JSONDecodeError:
            return cls(
                old_group_id="",
                root_pub_key="",
                curve=CurveType.SECP256K1,
                used_node_ids=[],
                old_threshold=0,
                new_threshold=0,
                new_node_ids=[],
                task_id="",
                biz_task_id="",
            )


# =================== extra_info ===================
@dataclass
class KeyGenRequestInfo:
    """Key generation request info"""

    org: Optional["cobo_waas2.OrgInfo"] = None
    project: Optional["cobo_waas2.MPCProject"] = None
    vault: Optional["cobo_waas2.MPCVault"] = None
    target_key_share_holder_group: Optional["cobo_waas2.KeyShareHolderGroup"] = None
    tss_request: Optional["cobo_waas2.TSSRequest"] = None

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeyGenRequestInfo instance from a JSON string.

        Args:
            json_str: JSON string containing the key generation request info

        Returns:
            KeyGenRequestInfo: An instance with parsed data or default values
        """
        if not json_str:
            return cls()

        try:
            data = json.loads(json_str)

            # Parse each nested object if present
            org = (
                cobo_waas2.OrgInfo.from_json(json.dumps(data.get("org")))
                if data.get("org")
                else None
            )
            project = (
                cobo_waas2.MPCProject.from_json(json.dumps(data.get("project")))
                if data.get("project")
                else None
            )
            vault = (
                cobo_waas2.MPCVault.from_json(json.dumps(data.get("vault")))
                if data.get("vault")
                else None
            )
            holder_group = (
                cobo_waas2.KeyShareHolderGroup.from_json(
                    json.dumps(data.get("target_key_share_holder_group"))
                )
                if data.get("target_key_share_holder_group")
                else None
            )
            tss_request = (
                cobo_waas2.TSSRequest.from_json(json.dumps(data.get("tss_request")))
                if data.get("tss_request")
                else None
            )

            return cls(
                org=org,
                project=project,
                vault=vault,
                target_key_share_holder_group=holder_group,
                tss_request=tss_request,
            )

        except json.JSONDecodeError:
            return cls()

    def __str__(self):
        if not self:
            return ""

        return json.dumps(
            {
                "org": self.org.to_dict() if self.org else None,
                "project": self.project.to_dict() if self.project else None,
                "vault": self.vault.to_dict() if self.vault else None,
                "target_key_share_holder_group": self.target_key_share_holder_group.to_dict()
                if self.target_key_share_holder_group
                else None,
                "tss_request": self.tss_request.to_dict() if self.tss_request else None,
            }
        )


@dataclass
class KeySignRequestInfo:
    """Key signing request info"""

    org: Optional[cobo_waas2.OrgInfo] = None
    project: Optional[cobo_waas2.MPCProject] = None
    vault: Optional[cobo_waas2.MPCVault] = None
    wallet: Optional[cobo_waas2.WalletInfo] = None
    signer_key_share_holder_group: Optional[cobo_waas2.KeyShareHolderGroup] = None
    source_addresses: Optional[List[cobo_waas2.AddressInfo]] = None
    transaction: Optional[cobo_waas2.Transaction] = None

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeySignRequestInfo instance from a JSON string.

        Args:
            json_str: JSON string containing the key signing request info

        Returns:
            KeySignRequestInfo: An instance with parsed data or default values
        """
        if not json_str:
            return cls()

        try:
            data = json.loads(json_str)

            # Parse each nested object if present
            org = (
                cobo_waas2.OrgInfo.from_json(json.dumps(data.get("org")))
                if data.get("org")
                else None
            )
            project = (
                cobo_waas2.MPCProject.from_json(json.dumps(data.get("project")))
                if data.get("project")
                else None
            )
            vault = (
                cobo_waas2.MPCVault.from_json(json.dumps(data.get("vault")))
                if data.get("vault")
                else None
            )
            wallet = (
                cobo_waas2.WalletInfo.from_json(json.dumps(data.get("wallet")))
                if data.get("wallet")
                else None
            )
            holder_group = (
                cobo_waas2.KeyShareHolderGroup.from_json(
                    json.dumps(data.get("signer_key_share_holder_group"))
                )
                if data.get("signer_key_share_holder_group")
                else None
            )

            # Handle list of addresses
            source_addresses = None
            if data.get("source_addresses"):
                source_addresses = [
                    cobo_waas2.AddressInfo.from_json(json.dumps(addr))
                    for addr in data.get("source_addresses")
                ]

            transaction = (
                cobo_waas2.Transaction.from_json(json.dumps(data.get("transaction")))
                if data.get("transaction")
                else None
            )

            return cls(
                org=org,
                project=project,
                vault=vault,
                wallet=wallet,
                signer_key_share_holder_group=holder_group,
                source_addresses=source_addresses,
                transaction=transaction,
            )

        except json.JSONDecodeError:
            return cls()

    def __str__(self):
        if not self:
            return ""
        return json.dumps(
            {
                "org": self.org.to_dict() if self.org else None,
                "project": self.project.to_dict() if self.project else None,
                "vault": self.vault.to_dict() if self.vault else None,
                "wallet": self.wallet.to_dict() if self.wallet else None,
                "signer_key_share_holder_group": self.signer_key_share_holder_group.to_dict()
                if self.signer_key_share_holder_group
                else None,
                "source_addresses": [addr.to_dict() for addr in self.source_addresses]
                if self.source_addresses
                else None,
                "transaction": self.transaction.to_dict() if self.transaction else None,
            }
        )


@dataclass
class KeyReshareRequestInfo:
    """Key reshare request info"""

    org: Optional[cobo_waas2.OrgInfo] = None
    project: Optional[cobo_waas2.MPCProject] = None
    vault: Optional[cobo_waas2.MPCVault] = None
    source_key_share_holder_group: Optional[cobo_waas2.KeyShareHolderGroup] = None
    target_key_share_holder_group: Optional[cobo_waas2.KeyShareHolderGroup] = None
    tss_request: Optional[cobo_waas2.TSSRequest] = None

    @classmethod
    def from_json(cls, json_str: str):
        """
        Create a KeyReshareRequestInfo instance from a JSON string.

        Args:
            json_str: JSON string containing the key reshare request info

        Returns:
            KeyReshareRequestInfo: An instance with parsed data or default values
        """
        if not json_str:
            return cls()

        try:
            data = json.loads(json_str)

            # Parse each nested object if present
            org = (
                cobo_waas2.OrgInfo.from_json(json.dumps(data.get("org")))
                if data.get("org")
                else None
            )
            project = (
                cobo_waas2.MPCProject.from_json(json.dumps(data.get("project")))
                if data.get("project")
                else None
            )
            vault = (
                cobo_waas2.MPCVault.from_json(json.dumps(data.get("vault")))
                if data.get("vault")
                else None
            )
            source_group = (
                cobo_waas2.KeyShareHolderGroup.from_json(
                    json.dumps(data.get("source_key_share_holder_group"))
                )
                if data.get("source_key_share_holder_group")
                else None
            )
            target_group = (
                cobo_waas2.KeyShareHolderGroup.from_json(
                    json.dumps(data.get("target_key_share_holder_group"))
                )
                if data.get("target_key_share_holder_group")
                else None
            )
            tss_request = (
                cobo_waas2.TSSRequest.from_json(json.dumps(data.get("tss_request")))
                if data.get("tss_request")
                else None
            )

            return cls(
                org=org,
                project=project,
                vault=vault,
                source_key_share_holder_group=source_group,
                target_key_share_holder_group=target_group,
                tss_request=tss_request,
            )

        except json.JSONDecodeError:
            return cls()

    def __str__(self):
        if not self:
            return ""
        return json.dumps(
            {
                "org": self.org.to_dict() if self.org else None,
                "project": self.project.to_dict() if self.project else None,
                "vault": self.vault.to_dict() if self.vault else None,
                "source_key_share_holder_group": self.source_key_share_holder_group.to_dict()
                if self.source_key_share_holder_group
                else None,
                "target_key_share_holder_group": self.target_key_share_holder_group.to_dict()
                if self.target_key_share_holder_group
                else None,
                "tss_request": self.tss_request.to_dict() if self.tss_request else None,
            }
        )
