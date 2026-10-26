<?php

declare(strict_types=1);

namespace WebFunctionClient;

use WebFunction\Page;
use WebFunction\Package;
use WebFunction\Pipeline;
use WebFunction\BadRequestError;
use WebFunction\UnexpectedStatusCodeError;
use WebFunction\JsonParseError;

// IMPORTANT - composer require webfunction-protocol/webfunction

/**
 * Generated client for the Generics Flawed API webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type ListOfAnyAttributes = array{total: int|float, items: list<mixed>}
 * @psalm-type ListOfAnyAttributes = array{total: int|float, items: list<mixed>}
 * @phpstan-type UserAttributes = array{id: string}
 * @psalm-type UserAttributes = array{id: string}
 * @phpstan-type ListOfGhostAttributes = array{total: int|float, items: list<array<string, mixed>>}
 * @psalm-type ListOfGhostAttributes = array{total: int|float, items: list<array<string, mixed>>}
 * @phpstan-type MixedOfStringArgs = array{plain?: string}
 * @psalm-type MixedOfStringArgs = array{plain?: string}
 * @phpstan-type ListOfUser2Attributes = array{total: int|float, items: list<UserAttributes>}
 * @psalm-type ListOfUser2Attributes = array{total: int|float, items: list<UserAttributes>}
 */
final class Client
{
    private \WebFunction\Client $client;

    public function __construct(?string $bearerAuth = null, ?string $version = null)
    {
        $this->client = \WebFunction\Client::fromPackageEndpoint('https://example.test/package', bearerAuth: $bearerAuth, version: $version);
    }

    /**
     * Escape hatch for calling an endpoint this client doesn't have a typed method for.
     */
    public function call(string $name, array $args = []): mixed
    {
        return $this->client->call($name, $args);
    }

    public function getPackage(): ?Package
    {
        return $this->client->getPackage();
    }

    public function setBearerAuth(?string $bearerAuth): void
    {
        $this->client->setBearerAuth($bearerAuth);
    }

    public function setVersion(?string $version): void
    {
        $this->client->setVersion($version);
    }

    public function setPipeline(?Pipeline $pipeline): void
    {
        $this->client->setPipeline($pipeline);
    }

    /**
     * @return ListOfAnyAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function unapplied(): array
    {
        return $this->client->call('unapplied');
    }

    /**
     * @return UserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function staticArg(): array
    {
        return $this->client->call('static-arg');
    }

    /**
     * @return mixed
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function emptyApp(): mixed
    {
        return $this->client->call('empty-app');
    }

    /**
     * @return mixed
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function badKey(): mixed
    {
        return $this->client->call('bad-key');
    }

    /**
     * @return mixed
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function badArg(): mixed
    {
        return $this->client->call('bad-arg');
    }

    /**
     * @return ListOfGhostAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function ghostArg(): array
    {
        return $this->client->call('ghost-arg');
    }

    /**
     * @param array{m?: MixedOfStringArgs} $args
     * @return UserAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function wrongContext(array $args = []): array
    {
        return $this->client->call('wrong-context', $args);
    }

    /**
     * @return ListOfUser2Attributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function collides(): array
    {
        return $this->client->call('collides');
    }

    /**
     * @return mixed
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function tReturn(): mixed
    {
        return $this->client->call('t-return');
    }

    /**
     * @param array{v?: mixed|string} $args
     * @return string
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function tArgument(array $args = []): string
    {
        return $this->client->call('t-argument', $args);
    }

    /**
     * @return ListOfAnyAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function tApplied(): array
    {
        return $this->client->call('t-applied');
    }

}

