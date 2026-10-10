<?php

declare(strict_types=1);

namespace OpenAllDebrid;

/**
 * A multipart/form-data body.
 *
 * @internal
 */
final class Multipart
{
    private readonly string $boundary;
    private string $body = '';

    public function __construct()
    {
        $this->boundary = 'open-alldebrid-' . bin2hex(random_bytes(12));
    }

    public function addFile(string $field, string $fileName, string $content, string $contentType): void
    {
        $name = addcslashes($field, "\"\\\r\n");
        $file = addcslashes($fileName, "\"\\\r\n");
        $this->body .= "--{$this->boundary}\r\n"
            . "Content-Disposition: form-data; name=\"{$name}\"; filename=\"{$file}\"\r\n"
            . "Content-Type: {$contentType}\r\n\r\n"
            . $content . "\r\n";
    }

    public function contentType(): string
    {
        return "multipart/form-data; boundary={$this->boundary}";
    }

    public function body(): string
    {
        return $this->body . "--{$this->boundary}--\r\n";
    }
}
