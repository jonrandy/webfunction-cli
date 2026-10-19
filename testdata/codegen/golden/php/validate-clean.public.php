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
 * Generated client for the Clean Example API webfunction package.
 * Source: https://example.test/package
 *
 * @phpstan-type PersonAttributes = array{id: string, name: string}
 * @psalm-type PersonAttributes = array{id: string, name: string}
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
     * @param array{id: string} $args
     * @return PersonAttributes
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function getPerson(array $args): array
    {
        return $this->client->call('get-person', $args);
    }

    /**
     * @param array{status?: 'active'|'inactive'} $args
     * @return ListPeoplePage
     * @throws UnexpectedStatusCodeError Non-200, non-400 response.
     * @throws JsonParseError Response body wasn't valid JSON.
     */
    public function listPeople(array $args = []): ListPeoplePage
    {
        $raw = $this->client->call('list-people', $args);
        /** @var Page $raw */
        return new ListPeoplePage($raw);
    }

}

/**
 * A page of results, wrapping the real \WebFunction\Page instance with
 * item-level typing on getItems() that the real class's own PHPDoc
 * (`@return array<int, mixed>`) can't provide by itself.
 */
final class ListPeoplePage implements \IteratorAggregate, \Countable
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

