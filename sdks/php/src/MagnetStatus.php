<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use OpenAllDebrid\Model\Magnet;

/** Magnet status codes. 5 and above are errors. */
final class MagnetStatus
{
    public const QUEUED = 0;
    public const DOWNLOADING = 1;
    public const COMPRESSING = 2;
    public const UPLOADING = 3;
    public const READY = 4;

    /** Whether the magnet's files are available. */
    public static function isReady(Magnet $magnet): bool
    {
        return $magnet->statusCode === self::READY;
    }

    /** Whether the magnet failed. */
    public static function isError(Magnet $magnet): bool
    {
        return ($magnet->statusCode ?? 0) > self::READY;
    }
}
