import React from 'react';

interface NotesIconProps extends React.ComponentPropsWithoutRef<'svg'> {
  size?: number;
}

export default function NotesIcon({ size, ...others }: NotesIconProps) {
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
      <rect x={5} y={3} width={14} height={18} rx={2} />
      <line x1={9} y1={7} x2={15} y2={7} />
      <line x1={9} y1={11} x2={15} y2={11} />
      <line x1={9} y1={15} x2={13} y2={15} />
    </svg>
  );
}
