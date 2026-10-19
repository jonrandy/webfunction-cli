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
 * Generated client for the Type Coverage Convert Test webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type AddressArgs = array{street: string, city?: string, postcode?: '10110'|'10120'}
 * @psalm-type AddressArgs = array{street: string, city?: string, postcode?: '10110'|'10120'}
 * @phpstan-type AddressAttributes = array{street: string, city: string, verified: bool}
 * @psalm-type AddressAttributes = array{street: string, city: string, verified: bool}
 * @phpstan-type AccountAttributes = array{id: string, email: string, age: int, active: bool, score: int|float, nickname: string|null, tags: list<string>, role: 'admin'|'member'|'guest', home: AddressAttributes, previous_homes: list<AddressAttributes>, extra: mixed}
 * @psalm-type AccountAttributes = array{id: string, email: string, age: int, active: bool, score: int|float, nickname: string|null, tags: list<string>, role: 'admin'|'member'|'guest', home: AddressAttributes, previous_homes: list<AddressAttributes>, extra: mixed}
 * @phpstan-type NodeAttributes = array{label: string, children: list<NodeAttributes>, parent: NodeAttributes|null}
 * @psalm-type NodeAttributes = array{label: string, children: list<NodeAttributes>, parent: NodeAttributes|null}
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
     * Exercises every argument type form.
     * @param array{email: string, id?: string, homepage?: string, feed?: string, birthday?: string, wake_time?: string, last_login?: string, phone?: string, avatar?: string, v4?: string, v6?: string, host?: string, age?: int, big?: int, delta?: int, huge_delta?: int, ratio?: float, precise?: float, joined_at?: int, score?: int|float, active?: bool, anything?: mixed, loose_list?: array, loose_object?: array<string, mixed>, tags?: list<string>, scores?: list<int|float>, nickname?: string|null, contact?: string|int|float, role?: 'admin'|'member'|'guest', level?: 1|2|3, home?: AddressArgs, previous_homes?: list<AddressArgs>} $args
     * @return AccountAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function createAccount(array $args): array
    {
        return $this->client->call('create-account', $args);
    }

    /**
     * @param array{limit?: int} $args
     * @return ListAccountsPage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listAccounts(array $args = []): ListAccountsPage
    {
        $raw = $this->client->call('list-accounts', $args);
        /** @var Page $raw */
        return new ListAccountsPage($raw);
    }

    /**
     * @return int|float
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getStats(): int|float
    {
        return $this->client->call('get-stats');
    }

    /**
     * @return bool
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function ping(): bool
    {
        return $this->client->call('ping');
    }

    /**
     * @return list<string>
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getTags(): array
    {
        return $this->client->call('get-tags');
    }

    /**
     * @param array{id: string} $args
     * @return AccountAttributes|null
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function findAccount(array $args): array|null
    {
        return $this->client->call('find-account', $args);
    }

    /**
     * @return NodeAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getTree(): array
    {
        return $this->client->call('get-tree');
    }

}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 */
final class ListAccountsPage implements \IteratorAggregate, \Countable
{
    private Page $page;

    public function __construct(Page $page)
    {
        $this->page = $page;
    }

    /**
     * @return list<mixed>
     */
    public function getItems(): array
    {
        return $this->page->getItems();
    }

    public function hasNext(): bool
    {
        return $this->page->hasNext();
    }

    public function hasPrevious(): bool
    {
        return $this->page->hasPrevious();
    }

    public function nextPage(): ?self
    {
        $next = $this->page->nextPage();
        return $next === null ? null : new self($next);
    }

    public function previousPage(): ?self
    {
        $prev = $this->page->previousPage();
        return $prev === null ? null : new self($prev);
    }

    public function count(): int
    {
        return $this->page->count();
    }

    public function getIterator(): \ArrayIterator
    {
        return $this->page->getIterator();
    }
}

