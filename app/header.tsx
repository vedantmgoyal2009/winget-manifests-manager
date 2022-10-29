import DarkModeIcon from './components/DarkModeIcon';
import LightModeIcon from './components/LightModeIcon';
import LoginIcon from './components/LoginIcon';
import LogoutIcon from './components/LogoutIcon';
import {
  ActionIcon,
  Group,
  Header,
  Switch,
  Text,
  useMantineColorScheme,
  useMantineTheme,
} from '@mantine/core';
import { auth as currentAuth, useAuth } from './auth';
import GithubIcon from './components/GitHubIcon';
import TwitterIcon from './components/TwitterIcon';

export default function AppHeader() {
  const auth = useAuth();
  const theme = useMantineTheme();
  const { colorScheme, toggleColorScheme } = useMantineColorScheme();

  return (
    <Header height={60}>
      <Group sx={{ height: '100%' }} px={20} position="apart">
        <Text
          variant="gradient"
          gradient={{ from: 'indigo', to: 'cyan', deg: 45 }}
          ta="center"
          fz="xl"
          fw={700}
        >
          {'Welcome to WinGet Community Manifests Manager!'}
        </Text>
        <Group spacing="lg">
          <ActionIcon
            variant="outline"
            radius="xl"
            size="lg"
            component="a"
            href="https://twitter.com/vedantmgoyal"
          >
            <TwitterIcon size={16} />
          </ActionIcon>
          <ActionIcon
            variant="outline"
            radius="xl"
            size="lg"
            component="a"
            href="https://github.com/vedantmgoyal2009/winget-manifests-manager"
          >
            <GithubIcon size={16} />
          </ActionIcon>
          <Switch
            checked={colorScheme === 'dark'}
            onChange={() => toggleColorScheme()}
            size="lg"
            onLabel={<LightModeIcon color={theme.white} size={20} />}
            offLabel={<DarkModeIcon color={theme.colors.gray[6]} size={20} />}
          />
          <ActionIcon
            variant="outline"
            onClick={async () => {
              currentAuth.currentUser !== null
                ? auth!.signout()
                : await auth!.signinWithGithub();
              window.location.reload();
            }}
            radius="xl"
            size="lg"
          >
            {currentAuth.currentUser !== null ? (
              <LogoutIcon size={16} />
            ) : (
              <LoginIcon size={16} />
            )}
          </ActionIcon>
        </Group>
      </Group>
    </Header>
  );
}
