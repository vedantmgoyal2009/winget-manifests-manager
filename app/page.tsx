'use client';
import { AppShell, Loader, useMantineColorScheme } from '@mantine/core';
import Editor from '@monaco-editor/react';
import { editor as monacoEditor } from 'monaco-editor/esm/vs/editor/editor.api';
import { useRef, useState } from 'react';
import AppHeader from './header';
import AppNavbar from './navbar';
import AppSidebar from './sidebar';

export default function Home() {
  const { colorScheme } = useMantineColorScheme();
  const [currentEditor, setCurrentEditor] = useState('Welcome');
  const editorRef = useRef<monacoEditor.IStandaloneCodeEditor | null>(null);
  let openEditors: {
    name: string;
    language: string;
    value: string;
    intialValue?: string;
  }[] = [
    {
      name: 'Welcome',
      language: 'javascript',
      value: '// Welcome to WinGet Community Manifests Manager...',
    },
    {
      name: 'File 1',
      language: 'yaml',
      value: 'Hello World:\n  - Hello World',
    },
  ];
  return (
    <AppShell
      padding="xl"
      layout="default"
      styles={(theme) => ({
        main: {
          backgroundColor:
            theme.colorScheme === 'dark'
              ? theme.colors.dark[8]
              : theme.colors.gray[0],
        },
      })}
      header={<AppHeader />}
      navbar={
        <AppNavbar
          openEditors={openEditors}
          currentEditor={currentEditor}
          setCurrentEditor={setCurrentEditor}
        />
      }
      aside={
        <AppSidebar
          openEditors={openEditors}
          setCurrentEditor={setCurrentEditor}
        />
      }
    >
      <Editor
        theme={colorScheme === 'dark' ? 'vs-dark' : 'vs-light'}
        loading={<Loader size="lg" variant="oval" />}
        language={
          openEditors.find((editor) => editor.name === currentEditor)!.language
        }
        value={
          openEditors.find((editor) => editor.name === currentEditor)!.value
        }
        onMount={(editor) => (editorRef.current = editor)}
        saveViewState={true}
      />
    </AppShell>
  );
}
