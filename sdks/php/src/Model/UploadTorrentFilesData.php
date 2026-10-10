<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class UploadTorrentFilesData
{
    public function __construct(
        /** @var list<UploadedMagnet> */
        public readonly array $files,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            files: Hydrate::list(Hydrate::required($data, 'files'), static fn (mixed $v): UploadedMagnet => UploadedMagnet::fromArray(Hydrate::object($v))),
        );
    }
}
