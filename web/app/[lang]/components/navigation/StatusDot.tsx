import React from 'react';

type Props = {
  connected: boolean;
};

const StatusDot = ({ connected }: Props) => {
  return (
    <svg width="20" height="10">
      <circle cx="5" cy="5" r="2.5" fill={connected ? 'lightgreen' : 'red'} />
    </svg>
  );
};

export default StatusDot;
