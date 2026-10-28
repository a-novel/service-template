import type { TemplateApi } from "./api";

import { HTTP_HEADERS } from "@a-novel-kit/nodelib-browser/http";

import { z } from "zod";

/** Parses an Item as returned by the service, coercing the timestamp strings to Date. */
export const ItemSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  description: z.string().optional(),
  createdAt: z.iso.datetime().transform((value) => new Date(value)),
  updatedAt: z.iso.datetime().transform((value) => new Date(value)),
});

/** An Item resource held by the template service. */
export type Item = z.infer<typeof ItemSchema>;

/** A page of Items, newest first, with how many exist in all. */
export const ItemPageSchema = z.object({
  items: z.array(ItemSchema),
  total: z.int(),
});

/** The value {@link ItemPageSchema} validates. */
export type ItemPage = z.infer<typeof ItemPageSchema>;

/** Creates a new Item and returns the stored record. */
export async function itemCreate(api: TemplateApi, name: string, description?: string): Promise<Item> {
  return await api.fetch("/v0/items", ItemSchema, {
    method: "POST",
    headers: HTTP_HEADERS.JSON,
    body: JSON.stringify({ name, description }),
  });
}

/** Fetches a single Item by ID. */
export async function itemGet(api: TemplateApi, id: string): Promise<Item> {
  const params = new URLSearchParams({ id });
  return await api.fetch(`/v0/item?${params.toString()}`, ItemSchema, { method: "GET", headers: HTTP_HEADERS.JSON });
}

/**
 * Lists a page of Items newest first, with their total, 100 at most when limit is omitted. Raise offset
 * by the page size to read the next page.
 */
export async function itemList(api: TemplateApi, limit?: number, offset?: number): Promise<ItemPage> {
  const params = new URLSearchParams({ limit: `${limit || 100}` });
  if (offset) params.set("offset", `${offset}`);
  return await api.fetch(`/v0/items?${params.toString()}`, ItemPageSchema, {
    method: "GET",
    headers: HTTP_HEADERS.JSON,
  });
}

/** Updates an existing Item and returns the stored record. */
export async function itemUpdate(api: TemplateApi, id: string, name: string, description?: string): Promise<Item> {
  return await api.fetch("/v0/item", ItemSchema, {
    method: "PUT",
    headers: HTTP_HEADERS.JSON,
    body: JSON.stringify({ id, name, description }),
  });
}

/** Deletes an Item by ID and returns the record that was removed. */
export async function itemDelete(api: TemplateApi, id: string): Promise<Item> {
  const params = new URLSearchParams({ id });
  return await api.fetch(`/v0/item?${params.toString()}`, ItemSchema, { method: "DELETE", headers: HTTP_HEADERS.JSON });
}
