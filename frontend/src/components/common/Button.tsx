import { type ButtonHTMLAttributes } from 'react'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary'
}

export default function Button({ 
  variant = 'primary', 
  className = '', 
  children, 
  ...props 
}: ButtonProps) {
  const baseClasses = "inline-flex items-center justify-center rounded-xl font-medium transition-all duration-200 px-6 py-3 cursor-pointer"
  
  const variantClasses = {
    primary: "bg-gray-900 hover:bg-gray-800 text-white active:scale-95",
    secondary: "bg-gray-100 hover:bg-gray-200 text-gray-700 active:scale-95"
  }

  return (
    <button 
      className={`${baseClasses} ${variantClasses[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  )
}
