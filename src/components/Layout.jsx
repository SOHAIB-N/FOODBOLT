import React from 'react';
import Logo from './Logo';

const Layout = ({ children }) => {
  return (
    <div className="layout">
      <header className="header">
        <Logo className="header-logo" />
        {/* Add navigation or other header elements here */}
      </header>
      <main className="main-content">
        {children}
      </main>
    </div>
  );
};

export default Layout; 