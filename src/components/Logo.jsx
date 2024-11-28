import React from 'react';

const Logo = ({ className = '' }) => {
  return (
    <div className={`logo-container ${className}`}>
      <img 
        src="/path/to/your/logo.png" 
        alt="Food Court Logo" 
        className="logo-image"
      />
    </div>
  );
};

export default Logo; 