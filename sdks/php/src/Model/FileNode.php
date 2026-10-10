<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

/** A file (`n`, `s`, `l`) or a folder (`n`, `e`) of a magnet. */
final class FileNode
{
    public function __construct(
        /** File or folder name. */
        public readonly string $n,
        /** Files only. Size in bytes. */
        public readonly ?int $s = null,
        /** Files only. Link to unlock with `/v4/link/unlock`. */
        public readonly ?string $l = null,
        /**
         * Folders only. The entries of the folder.
         * @var list<FileNode>|null
         */
        public readonly ?array $e = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            n: Hydrate::string(Hydrate::required($data, 'n')),
            s: array_key_exists('s', $data) && $data['s'] !== null ? Hydrate::int($data['s']) : null,
            l: array_key_exists('l', $data) && $data['l'] !== null ? Hydrate::string($data['l']) : null,
            e: array_key_exists('e', $data) && $data['e'] !== null ? Hydrate::list($data['e'], static fn (mixed $v): FileNode => FileNode::fromArray(Hydrate::object($v))) : null,
        );
    }
}
