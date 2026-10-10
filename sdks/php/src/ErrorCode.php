<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid;

/** Error codes the API is known to return. The API may return others. */
final class ErrorCode
{
    public const GENERIC = 'GENERIC';
    public const NOT_FOUND = '404';
    public const MAINTENANCE = 'MAINTENANCE';
    public const DISCONTINUED = 'DISCONTINUED';
    public const AUTH_MISSING_APIKEY = 'AUTH_MISSING_APIKEY';
    public const AUTH_BAD_APIKEY = 'AUTH_BAD_APIKEY';
    public const AUTH_BLOCKED = 'AUTH_BLOCKED';
    public const AUTH_USER_BANNED = 'AUTH_USER_BANNED';
    public const ALREADY_SENT = 'ALREADY_SENT';
    public const NO_SERVER = 'NO_SERVER';
    public const LINK_IS_MISSING = 'LINK_IS_MISSING';
    public const BAD_LINK = 'BAD_LINK';
    public const LINK_HOST_NOT_SUPPORTED = 'LINK_HOST_NOT_SUPPORTED';
    public const LINK_DOWN = 'LINK_DOWN';
    public const LINK_PASS_PROTECTED = 'LINK_PASS_PROTECTED';
    public const LINK_HOST_UNAVAILABLE = 'LINK_HOST_UNAVAILABLE';
    public const LINK_TOO_MANY_DOWNLOADS = 'LINK_TOO_MANY_DOWNLOADS';
    public const LINK_HOST_FULL = 'LINK_HOST_FULL';
    public const LINK_HOST_LIMIT_REACHED = 'LINK_HOST_LIMIT_REACHED';
    public const LINK_ERROR = 'LINK_ERROR';
    public const LINK_TEMPORARY_UNAVAILABLE = 'LINK_TEMPORARY_UNAVAILABLE';
    public const LINK_NOT_SUPPORTED = 'LINK_NOT_SUPPORTED';
    public const REDIRECTOR_NOT_SUPPORTED = 'REDIRECTOR_NOT_SUPPORTED';
    public const REDIRECTOR_ERROR = 'REDIRECTOR_ERROR';
    public const REDIRECTOR_RATE_LIMITED = 'REDIRECTOR_RATE_LIMITED';
    public const STREAM_INVALID_GEN_ID = 'STREAM_INVALID_GEN_ID';
    public const STREAM_INVALID_STREAM_ID = 'STREAM_INVALID_STREAM_ID';
    public const DELAYED_INVALID_ID = 'DELAYED_INVALID_ID';
    public const FREE_TRIAL_LIMIT_REACHED = 'FREE_TRIAL_LIMIT_REACHED';
    public const MUST_BE_PREMIUM = 'MUST_BE_PREMIUM';
    public const MAGNET_INVALID_ID = 'MAGNET_INVALID_ID';
    public const MAGNET_INVALID_URI = 'MAGNET_INVALID_URI';
    public const MAGNET_INVALID_FILE = 'MAGNET_INVALID_FILE';
    public const MAGNET_FILE_UPLOAD_FAILED = 'MAGNET_FILE_UPLOAD_FAILED';
    public const MAGNET_NO_URI = 'MAGNET_NO_URI';
    public const MAGNET_PROCESSING = 'MAGNET_PROCESSING';
    public const MAGNET_TOO_MANY_ACTIVE = 'MAGNET_TOO_MANY_ACTIVE';
    public const MAGNET_TOO_MANY = 'MAGNET_TOO_MANY';
    public const MAGNET_MUST_BE_PREMIUM = 'MAGNET_MUST_BE_PREMIUM';
    public const MAGNET_NO_SERVER = 'MAGNET_NO_SERVER';
    public const MAGNET_TOO_LARGE = 'MAGNET_TOO_LARGE';
    public const MAGNET_UPLOAD_FAILED = 'MAGNET_UPLOAD_FAILED';
    public const MAGNET_INTERNAL_ERROR = 'MAGNET_INTERNAL_ERROR';
    public const MAGNET_CANT_BOOTSTRAP = 'MAGNET_CANT_BOOTSTRAP';
    public const MAGNET_MAGNET_TOO_BIG = 'MAGNET_MAGNET_TOO_BIG';
    public const MAGNET_TOOK_TOO_LONG = 'MAGNET_TOOK_TOO_LONG';
    public const MAGNET_LINKS_REMOVED = 'MAGNET_LINKS_REMOVED';
    public const MAGNET_PROCESSING_FAILED = 'MAGNET_PROCESSING_FAILED';
    public const PIN_ALREADY_AUTHED = 'PIN_ALREADY_AUTHED';
    public const PIN_EXPIRED = 'PIN_EXPIRED';
    public const PIN_INVALID = 'PIN_INVALID';
    public const USER_LINK_MISSING = 'USER_LINK_MISSING';
    public const USER_LINK_INVALID = 'USER_LINK_INVALID';
    public const MISSING_NOTIF_ENDPOINT = 'MISSING_NOTIF_ENDPOINT';
    public const VOUCHER_DURATION_INVALID = 'VOUCHER_DURATION_INVALID';
    public const VOUCHER_NB_INVALID = 'VOUCHER_NB_INVALID';
    public const NO_MORE_VOUCHER = 'NO_MORE_VOUCHER';
    public const INSUFFICIENT_BALANCE = 'INSUFFICIENT_BALANCE';
    public const DOWNLOAD_FAILED = 'DOWNLOAD_FAILED';
    public const ACCOUNT_INVALID = 'ACCOUNT_INVALID';
    public const NO_JSON_PARAM = 'NO_JSON_PARAM';
    public const JSON_INVALID = 'JSON_INVALID';
    public const FREEDAYS_INVALID_COUNTRY = 'FREEDAYS_INVALID_COUNTRY';
    public const FREEDAYS_INVALID_PHONE = 'FREEDAYS_INVALID_PHONE';
    public const FREEDAYS_INVALID_PROVIDER = 'FREEDAYS_INVALID_PROVIDER';
    public const FREEDAYS_USED_PHONE = 'FREEDAYS_USED_PHONE';
    public const FREEDAYS_ALREADY_REQUESTED = 'FREEDAYS_ALREADY_REQUESTED';
    public const FREEDAYS_INVALID_STATUS = 'FREEDAYS_INVALID_STATUS';
    public const FREEDAYS_TOO_MUCH_RETRIES = 'FREEDAYS_TOO_MUCH_RETRIES';
    public const FREEDAYS_INVALID_CODE = 'FREEDAYS_INVALID_CODE';
    public const FREEDAYS_DENIED = 'FREEDAYS_DENIED';
    public const FREEDAYS_ERROR_SENDING = 'FREEDAYS_ERROR_SENDING';

    /** Every code of the spec. */
    public const ALL = [
        self::GENERIC,
        self::NOT_FOUND,
        self::MAINTENANCE,
        self::DISCONTINUED,
        self::AUTH_MISSING_APIKEY,
        self::AUTH_BAD_APIKEY,
        self::AUTH_BLOCKED,
        self::AUTH_USER_BANNED,
        self::ALREADY_SENT,
        self::NO_SERVER,
        self::LINK_IS_MISSING,
        self::BAD_LINK,
        self::LINK_HOST_NOT_SUPPORTED,
        self::LINK_DOWN,
        self::LINK_PASS_PROTECTED,
        self::LINK_HOST_UNAVAILABLE,
        self::LINK_TOO_MANY_DOWNLOADS,
        self::LINK_HOST_FULL,
        self::LINK_HOST_LIMIT_REACHED,
        self::LINK_ERROR,
        self::LINK_TEMPORARY_UNAVAILABLE,
        self::LINK_NOT_SUPPORTED,
        self::REDIRECTOR_NOT_SUPPORTED,
        self::REDIRECTOR_ERROR,
        self::REDIRECTOR_RATE_LIMITED,
        self::STREAM_INVALID_GEN_ID,
        self::STREAM_INVALID_STREAM_ID,
        self::DELAYED_INVALID_ID,
        self::FREE_TRIAL_LIMIT_REACHED,
        self::MUST_BE_PREMIUM,
        self::MAGNET_INVALID_ID,
        self::MAGNET_INVALID_URI,
        self::MAGNET_INVALID_FILE,
        self::MAGNET_FILE_UPLOAD_FAILED,
        self::MAGNET_NO_URI,
        self::MAGNET_PROCESSING,
        self::MAGNET_TOO_MANY_ACTIVE,
        self::MAGNET_TOO_MANY,
        self::MAGNET_MUST_BE_PREMIUM,
        self::MAGNET_NO_SERVER,
        self::MAGNET_TOO_LARGE,
        self::MAGNET_UPLOAD_FAILED,
        self::MAGNET_INTERNAL_ERROR,
        self::MAGNET_CANT_BOOTSTRAP,
        self::MAGNET_MAGNET_TOO_BIG,
        self::MAGNET_TOOK_TOO_LONG,
        self::MAGNET_LINKS_REMOVED,
        self::MAGNET_PROCESSING_FAILED,
        self::PIN_ALREADY_AUTHED,
        self::PIN_EXPIRED,
        self::PIN_INVALID,
        self::USER_LINK_MISSING,
        self::USER_LINK_INVALID,
        self::MISSING_NOTIF_ENDPOINT,
        self::VOUCHER_DURATION_INVALID,
        self::VOUCHER_NB_INVALID,
        self::NO_MORE_VOUCHER,
        self::INSUFFICIENT_BALANCE,
        self::DOWNLOAD_FAILED,
        self::ACCOUNT_INVALID,
        self::NO_JSON_PARAM,
        self::JSON_INVALID,
        self::FREEDAYS_INVALID_COUNTRY,
        self::FREEDAYS_INVALID_PHONE,
        self::FREEDAYS_INVALID_PROVIDER,
        self::FREEDAYS_USED_PHONE,
        self::FREEDAYS_ALREADY_REQUESTED,
        self::FREEDAYS_INVALID_STATUS,
        self::FREEDAYS_TOO_MUCH_RETRIES,
        self::FREEDAYS_INVALID_CODE,
        self::FREEDAYS_DENIED,
        self::FREEDAYS_ERROR_SENDING,
    ];
}
