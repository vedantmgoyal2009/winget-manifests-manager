import {
  Aside,
  Box,
  Button,
  Divider,
  Group,
  Select,
  Stack,
  Tabs,
  Textarea,
} from '@mantine/core';
import { Dispatch, SetStateAction, useState } from 'react';
import { firestore } from './auth';
import { collection, getDocs } from 'firebase/firestore/lite';
import useSWR from 'swr';

export default function AppSidebar({
  openEditors,
  setCurrentEditor,
}: {
  openEditors: {
    name: string;
    language: string;
    value: string;
    initialValue?: string;
  }[];
  setCurrentEditor: Dispatch<SetStateAction<string>>;
}) {
  const req1 = useSWR<Record<string, string[]>>(
    '/api/get-winget-packages',
    async (url: string) => {
      const res = await fetch(url);
      return res.json();
    },
  );
  const req2 = useSWR<string[]>(
    'wpa-packages',
    async (collectionName: string) =>
      (await getDocs(collection(firestore, collectionName))).docs.map(
        (doc) => doc.id,
      ),
  );
  const [pkgVersions, automationPkgs] = [req1.data, req2.data];
  const openNewEditor = (
    name: string,
    language: string,
    value: string,
  ): void => {
    openEditors.push({ name, language, value, initialValue: value });
    setCurrentEditor(name);
  };
  const getEditor = (editorName: string): string =>
    openEditors.find((editor) => editor.name === editorName)!.value;
  const hasEditorBeenModified = (editorName: string): boolean => {
    const editor = openEditors.find((editor) => editor.name === editorName);
    return editor!.value !== editor!.initialValue;
  };
  const [selectedPkgId, setPkgId] = useState<string>('');
  const [selectedPkgVersion, setPkgVersion] = useState<string>('');
  const [isNewPkg, setIsNewPkg] = useState<boolean>(false);
  const [testResult, setTestResult] = useState<string>('To be implemented...');
  return (
    <Aside width={{ base: 303 }}>
      <Aside.Section grow>
        <Tabs defaultValue="automation-manager" variant="default" radius="sm">
          <Tabs.List position="center">
            <Tabs.Tab value="automation-manager">Automation</Tabs.Tab>
            <Tabs.Tab value="create-locales">Locales</Tabs.Tab>
            <Tabs.Tab value="edit-manifest">Edit Manifests</Tabs.Tab>
          </Tabs.List>
          <Tabs.Panel value="automation-manager">
            <Box
              component={Stack}
              align="stretch"
              justify="flex-end"
              spacing="xs"
              sx={(theme) => ({
                padding: theme.spacing.sm,
                paddingTop: theme.spacing.sm / 7,
              })}
            >
              <Select
                searchable
                clearable
                required
                creatable
                label="Select package identifier"
                placeholder="Pick or create a package..."
                data={automationPkgs || []}
                onChange={(value) => {
                  setIsNewPkg(!automationPkgs!.includes(value as string));
                  setPkgId(value || '');
                }}
                getCreateLabel={(value) => `+ Add ${value}`}
                onCreate={(value) => value}
              />
              <Group grow position="apart" spacing="xs">
                <Button
                  variant="light"
                  radius="xl"
                  disabled={!isNewPkg || selectedPkgId === ''}
                  children="Add Package"
                />
                <Button
                  variant="outline"
                  radius="xl"
                  disabled={isNewPkg || selectedPkgId === ''}
                  children="Modify"
                />
              </Group>
              <Button
                variant="light"
                radius="xl"
                color="teal"
                children="Test Package for Automation"
              />
              <Button
                variant="outline"
                radius="xl"
                color="orange"
                children="Pause/Resume Updates for Package"
              />
              <Button
                variant="outline"
                radius="xl"
                color="red"
                children="Remove Package (Cannot be undone)"
              />
            </Box>
          </Tabs.Panel>
          <Tabs.Panel value="create-locales">x</Tabs.Panel>
          <Tabs.Panel value="edit-manifest">y</Tabs.Panel>
        </Tabs>
      </Aside.Section>
      <Aside.Section
        sx={(theme) => ({
          padding: theme.spacing.sm,
        })}
      >
        <Divider my="sm" />
        <Textarea disabled minRows={8} label="Test Result" value={testResult} />
      </Aside.Section>
    </Aside>
  );
}
