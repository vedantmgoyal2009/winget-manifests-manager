import React from 'react';

interface LightModeIconProps extends React.ComponentPropsWithoutRef<'svg'> {
  size?: number;
}

export default function LightModeIcon({ size, ...others }: LightModeIconProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      strokeWidth={2}
      stroke="currentColor"
      strokeLinecap="round"
      strokeLinejoin="round"
      width={size}
      height={size}
      {...others}
    >
      <path stroke="none" d="M0 0h24v24H0z" fill="none" />
      <path d="M3 13h1" />
      <path d="M20 13h1" />
      <path d="M5.6 6.6l.7 .7" />
      <path d="M18.4 6.6l-.7 .7" />
      <path d="M8 13a4 4 0 1 1 8 0" />
      <path d="M3 17h18" />
      <path d="M7 20h5" />
      <path d="M16 20h1" />
      <path d="M12 5v-1" />
    </svg>
  );
}
