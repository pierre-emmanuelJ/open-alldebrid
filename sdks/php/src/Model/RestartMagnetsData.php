<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class RestartMagnetsData
{
    public function __construct(
        /** With `id`. The restarted magnet. */
        public readonly ?int $id = null,
        /** With `id`. A confirmation message. */
        public readonly ?string $message = null,
        /**
         * With `ids[]`. One entry per magnet.
         * @var list<RestartedMagnet>|null
         */
        public readonly ?array $magnets = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            id: array_key_exists('id', $data) && $data['id'] !== null ? Hydrate::int($data['id']) : null,
            message: array_key_exists('message', $data) && $data['message'] !== null ? Hydrate::string($data['message']) : null,
            magnets: array_key_exists('magnets', $data) && $data['magnets'] !== null ? Hydrate::list($data['magnets'], static fn (mixed $v): RestartedMagnet => RestartedMagnet::fromArray(Hydrate::object($v))) : null,
        );
    }
}
