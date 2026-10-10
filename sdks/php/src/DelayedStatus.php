<?php

declare(strict_types=1);

namespace OpenAllDebrid;

/** Delayed link statuses. */
final class DelayedStatus
{
    public const PROCESSING = 1;
    public const READY = 2;
    public const FAILED = 3;
}
