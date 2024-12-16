import React from 'react';
import SessionComponent from '../components/Session';
import { useLocation } from 'react-router-dom';

const useQuery = () => {
  return new URLSearchParams(useLocation().search);
};

const SessionPage = () => {
  const query = useQuery();
  const sessionId = query.get('session');
  const userId = localStorage.getItem('User-ID');

  return (
    <div>
      <SessionComponent sessionId={sessionId} userId={userId} />
    </div>
  );
};

export default SessionPage;