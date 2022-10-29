'use client';
import { useState } from 'react';
import { CacheProvider } from '@emotion/react';
import {
  ColorScheme,
  ColorSchemeProvider,
  MantineProvider,
  useEmotionCache,
} from '@mantine/core';
import { useServerInsertedHTML } from 'next/navigation';
import { AuthProvider } from './auth';

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const cache = useEmotionCache();
  cache.compat = true;
  useServerInsertedHTML(() => {
    return (
      <style
        data-emotion={`${cache.key} ${Object.keys(cache.inserted).join(' ')}`}
        dangerouslySetInnerHTML={{
          __html: Object.values(cache.inserted).join(' '),
        }}
      />
    );
  });
  const [colorScheme, setColorScheme] = useState<ColorScheme>('light');
  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));
  return (
    <html lang="en">
      {/*
        <head /> will contain the components returned by the nearest parent
        head.tsx. Find out more at https://beta.nextjs.org/docs/api-reference/file-conventions/head
      */}
      {/*<head />*/}
      <head>
        <title>WinGet Community Manifests Manager</title>
        <meta content="width=device-width, initial-scale=1" name="viewport" />
        <meta
          name="description"
          content="By the community, for the community."
        />
        <link rel="icon" href="/favicon.ico" />
      </head>
      <body>
        <CacheProvider value={cache}>
          <ColorSchemeProvider
            colorScheme={colorScheme}
            toggleColorScheme={toggleColorScheme}
          >
            <MantineProvider
              withGlobalStyles
              withNormalizeCSS
              theme={{ colorScheme }}
            >
              <AuthProvider>{children}</AuthProvider>
            </MantineProvider>
          </ColorSchemeProvider>
        </CacheProvider>
      </body>
    </html>
  );
}
