<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class MagnetFiles
{
    public function __construct(
        /** The magnet ID (sent back as a string). */
        public readonly string $id,
        /** @var list<FileNode>|null */
        public readonly ?array $files = null,
        public readonly ?ApiError $error = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            id: Hydrate::flexString(Hydrate::required($data, 'id')),
            files: array_key_exists('files', $data) && $data['files'] !== null ? Hydrate::list($data['files'], static fn (mixed $v): FileNode => FileNode::fromArray(Hydrate::object($v))) : null,
            error: array_key_exists('error', $data) && $data['error'] !== null ? ApiError::fromArray(Hydrate::object($data['error'])) : null,
        );
    }
}
