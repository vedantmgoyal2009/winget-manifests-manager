// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import { initializeApp, cert } from 'firebase-admin/app';
import { getAuth } from 'firebase-admin/auth';
import type { RequestOptions } from 'http';
import type { NextApiRequest, NextApiResponse } from 'next';

type Response = {
  text?: string;
  error?: string;
};

const app = initializeApp({
  credential: cert({
    clientEmail: process.env.FIREBASE_CLIENT_EMAIL,
    privateKey: process.env.FIREBASE_PRIVATE_KEY,
    projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  }),
  databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
});

export default async (req: NextApiRequest, res: NextApiResponse<Response>) => {
  if (!req.headers.authorization) {
    return res
      .status(401)
      .json({ error: 'The API is not meant to be called directly.' });
  }

  try {
    await getAuth(app).verifyIdToken(req.headers.authorization.split(' ')[1]);
  } catch (error) {
    res.status(500).json({ error: 'Authentication failed.' });
  }

  await translate(req.body)
    .then((response) => res.status(200).json(response))
    .catch((error) => res.status(500).json({ error }));
};

async function translate(inputText: string, options?: TranslateOptions) {
  return new Translator(inputText, options).translate();
}

class Translator {
  protected options: typeof defaults & TranslateOptions;

  constructor(protected inputText: string, options?: TranslateOptions) {
    this.options = Object.assign({}, defaults, options);
  }

  async translate() {
    const url = this.buildUrl();
    const fetchOptions = this.buildFetchOptions();
    const res = await fetch(url, fetchOptions);
    if (!res.ok)
      throw new Error(`Translate API returned an error: ${res.statusText}`);
    const raw = (await res.json()) as RawResponse;
    const text = this.buildResText(raw);
    return { text, raw };
  }

  protected buildUrl() {
    const { host } = this.options;
    return [
      `https://${host}/translate_a/single`,
      '?client=at',
      '&dt=t', // return sentences
      '&dt=rm', // add translit to sentences
      '&dj=1', // result as pretty json instead of deep nested arrays
    ].join('');
  }

  protected buildBody() {
    const { from, to } = this.options;
    const params = {
      sl: from,
      tl: to,
      q: this.inputText,
    };
    return new URLSearchParams(params).toString();
  }

  protected buildFetchOptions() {
    const { fetchOptions } = this.options;
    const res = Object.assign({}, fetchOptions);
    res.method = 'POST';
    res.headers = Object.assign({}, res.headers, {
      'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8',
    });
    res.body = this.buildBody();
    return res;
  }

  protected buildResText({ sentences }: RawResponse) {
    return sentences
      .filter((s): s is Sentence => 'trans' in s)
      .map((s) => s.trans)
      .join('');
  }
}

const defaults: Required<Pick<TranslateOptions, 'from' | 'to' | 'host'>> = {
  from: 'auto',
  to: 'en',
  host: 'translate.google.com',
};

interface TranslateOptions {
  from?: string;
  to?: string;
  host?: string;
  fetchOptions?: Partial<RequestInit>;
}

interface RawResponse {
  sentences: (Sentence | SrcTranslit)[];
  src: string;
  confidence: number;
  ld_result: {
    srclangs: string[];
    srclangs_confidences: number[];
    extended_srclangs: string[];
  };
}

interface Sentence {
  trans: string;
  orig: string;
}

interface SrcTranslit {
  src_translit: string;
}

interface RequestInit {
  /**
   * A BodyInit object or null to set request's body.
   */
  body?: BodyInit | null;
  /**
   * A Headers object, an object literal, or an array of two-item arrays to set request's headers.
   */
  headers?: HeadersInit;
  /**
   * A string to set request's method.
   */
  method?: string;
  /**
   * A string indicating whether request follows redirects, results in an error upon encountering a redirect, or returns the redirect (in an opaque fashion). Sets request's redirect.
   */
  redirect?: RequestRedirect;
  /**
   * An AbortSignal to set request's signal.
   */
  signal?: AbortSignal | null;
  /**
   * A string whose value is a same-origin URL, "about:client", or the empty string, to set request’s referrer.
   */
  referrer?: string;
  /**
   * A referrer policy to set request’s referrerPolicy.
   */
  referrerPolicy?: ReferrerPolicy;

  // Node-fetch extensions to the whatwg/fetch spec
  agent?:
    | RequestOptions['agent']
    | ((parsedUrl: URL) => RequestOptions['agent']);
  compress?: boolean;
  counter?: number;
  follow?: number;
  hostname?: string;
  port?: number;
  protocol?: string;
  size?: number;
  highWaterMark?: number;
  insecureHTTPParser?: boolean;
}
