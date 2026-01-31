import { type InputHTMLAttributes } from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
}

export default function Input({ label, className = '', ...props }: InputProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-bold text-slate-700 mb-2">
          {label}
        </label>
      )}
      <input
        className={`
          flex h-11 w-full rounded-2xl border-none 
          bg-slate-100 px-4 py-2 text-sm text-slate-900 
          placeholder:text-slate-400 
          focus:bg-white focus:ring-2 focus:ring-blue-500/20 
          transition-all duration-200 outline-none
          ${className}
        `}
        {...props}
      />
    </div>
  )
}
