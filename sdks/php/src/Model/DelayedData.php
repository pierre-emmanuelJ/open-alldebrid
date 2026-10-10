<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class DelayedData
{
    public function __construct(
        /** `1` still processing, `2` ready (`link` is set), `3` failed. */
        public readonly int $status,
        /** Estimated seconds left. */
        public readonly int $timeLeft,
        /** The download link, once `status` is `2`. */
        public readonly ?string $link = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            status: Hydrate::int(Hydrate::required($data, 'status')),
            timeLeft: Hydrate::int(Hydrate::required($data, 'time_left')),
            link: array_key_exists('link', $data) && $data['link'] !== null ? Hydrate::string($data['link']) : null,
        );
    }
}
