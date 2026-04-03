import { HttpClient } from '../httpClient';

describe('HttpClient', () => {
  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('treats +json media types as JSON responses', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ ok: true }), {
        headers: {
          'content-type': 'application/vnd.api+json'
        },
        status: 200
      })
    );

    vi.stubGlobal('fetch', fetchMock);

    const client = new HttpClient('/api');

    await expect(client.get<{ ok: boolean }>('/rooms')).resolves.toEqual({ ok: true });
  });

  it('returns an ApiError fallback when a +json error body is invalid JSON', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response('Upstream gateway failure', {
        headers: {
          'content-type': 'application/problem+json'
        },
        status: 502
      })
    );

    vi.stubGlobal('fetch', fetchMock);

    const client = new HttpClient('/api');

    await expect(client.get('/rooms')).rejects.toMatchObject({
      message: 'Upstream gateway failure',
      status: 502
    });
  });
});
