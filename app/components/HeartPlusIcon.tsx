import React from 'react';

interface HeartPlusIconProps extends React.ComponentPropsWithoutRef<'svg'> {
  size?: number;
}

export default function HeartPlusIcon({ size, ...others }: HeartPlusIconProps) {
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
      <path d="M13 19l-1 1l-7.5 -7.428a5 5 0 1 1 7.5 -6.566a5 5 0 0 1 8.003 5.996" />
      <path d="M14 16h6" />
      <path d="M17 13v6" />
    </svg>
  );
}
