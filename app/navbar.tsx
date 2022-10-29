import {
  Avatar,
  Box,
  Group,
  Navbar,
  Text,
  ThemeIcon,
  UnstyledButton,
  createStyles,
  NavLink,
} from '@mantine/core';
import { Dispatch, SetStateAction } from 'react';
import { auth as currentAuth } from './auth';
import NotebookIcon from './components/NotebookIcon';
import NotesIcon from './components/NotesIcon';

export default function AppNavbar({
  openEditors,
  currentEditor,
  setCurrentEditor,
}: {
  openEditors: {
    name: string;
    language: string;
    value: string;
    initialValue?: string;
  }[];
  currentEditor: string;
  setCurrentEditor: Dispatch<SetStateAction<string>>;
}) {
  const { classes, theme } = createStyles((theme) => ({
    navLink: {
      fontWeight: 500,
      textDecoration: 'none',
      padding: `${theme.spacing.xs}px ${theme.spacing.md}px`,
      fontSize: theme.fontSizes.sm,
      color:
        theme.colorScheme === 'dark'
          ? theme.colors.dark[0]
          : theme.colors.gray[7],
      '&:hover': {
        backgroundColor:
          theme.colorScheme === 'dark'
            ? theme.colors.dark[7]
            : theme.colors.gray[0],
        color: theme.colorScheme === 'dark' ? theme.white : theme.black,
      },
    },
  }))();
  return (
    <Navbar p="xs" width={{ base: 303 }}>
      <Navbar.Section grow mt="md">
        <NavLink
          defaultOpened
          label="Opened Editors"
          icon={
            <ThemeIcon variant="light" size={30}>
              <NotebookIcon size={18} />
            </ThemeIcon>
          }
          childrenOffset={28}
          className={classes.navLink}
        >
          {openEditors.map((editor) => {
            return (
              <NavLink
                active={currentEditor === editor.name}
                icon={<NotesIcon size={18} />}
                id="editor-inside-opened-editors"
                label={editor.name}
                className={classes.navLink}
                style={{
                  borderLeft: `1px solid ${
                    theme.colorScheme === 'dark'
                      ? theme.colors.dark[4]
                      : theme.colors.gray[3]
                  }`,
                }}
                onClick={() => setCurrentEditor(editor.name)}
              />
            );
          })}
        </NavLink>
      </Navbar.Section>
      <Navbar.Section>
        <Box
          sx={{
            paddingTop: theme.spacing.sm,
            borderTop: `1px solid ${
              theme.colorScheme === 'dark'
                ? theme.colors.dark[4]
                : theme.colors.gray[2]
            }`,
          }}
        >
          <UnstyledButton
            sx={{
              display: 'block',
              width: '100%',
              padding: theme.spacing.xs,
              borderRadius: theme.radius.sm,
              color:
                theme.colorScheme === 'dark'
                  ? theme.colors.dark[0]
                  : theme.black,
              backgroundColor:
                theme.colorScheme === 'dark'
                  ? theme.colors.dark[6]
                  : theme.colors.gray[1],
              cursor: 'not-allowed',
            }}
          >
            <Group>
              <Avatar
                size="lg"
                variant="filled"
                radius="xl"
                src={currentAuth.currentUser?.photoURL}
              />
              <Box sx={{ flex: 1 }}>
                <Text size="sm" weight={500}>
                  {currentAuth.currentUser?.displayName || 'Unauthorized'}
                </Text>
                <Text size="xs" color="dimmed">
                  {currentAuth.currentUser?.email || 'Waiting for login...'}
                </Text>
              </Box>
            </Group>
          </UnstyledButton>
        </Box>
      </Navbar.Section>
    </Navbar>
  );
}
