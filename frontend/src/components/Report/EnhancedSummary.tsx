import React, { type ReactNode } from 'react';
import ReactMarkdown from 'react-markdown';
import { Bot, Sparkles, AlertTriangle, Target, CheckCircle2 } from 'lucide-react';
import type { ReportData, BrandRanking, BrandAnalysis } from '../../types/report';

interface EnhancedSummaryProps {
  data?: ReportData;
  recommendation?: string;
  rankings?: BrandRanking[];
  brandAnalysis?: Record<string, BrandAnalysis>;
}

interface MarkdownComponentProps {
  children?: ReactNode
}

/**
 * EnhancedSummary Component
 * 
 * Displays the AI-generated purchase advice, scenario-based recommendations,
 * and avoidance guide (pitfalls) using the report recommendation data.
 */
export const EnhancedSummary: React.FC<EnhancedSummaryProps> = (props) => {
  const recommendation = props.recommendation || props.data?.recommendation;
  const rankingsList = props.rankings || props.data?.rankings;

  if (!recommendation) {
    return null;
  }

  // Custom components for ReactMarkdown to enhance styling
  const markdownComponents = {
    h1: ({ children }: MarkdownComponentProps) => (
      <h2 className="text-xl font-bold text-gray-900 mt-6 mb-4 flex items-center gap-2">{children}</h2>
    ),
    h2: ({ children }: MarkdownComponentProps) => (
      <h3 className="text-lg font-bold text-gray-800 mt-5 mb-3 flex items-center gap-2 border-l-4 border-blue-500 pl-3">{children}</h3>
    ),
    h3: ({ children }: MarkdownComponentProps) => (
      <h4 className="text-base font-bold text-gray-700 mt-4 mb-2">{children}</h4>
    ),
    ul: ({ children }: MarkdownComponentProps) => (
      <ul className="space-y-2 mb-4 text-gray-600 list-none">{children}</ul>
    ),
    ol: ({ children }: MarkdownComponentProps) => (
      <ol className="space-y-2 mb-4 text-gray-600 list-decimal list-inside">{children}</ol>
    ),
    li: ({ children }: MarkdownComponentProps) => (
      <li className="flex items-start gap-2">
        <span className="mt-1.5 w-1.5 h-1.5 rounded-full bg-blue-400 flex-shrink-0 block" />
        <span className="flex-1">{children}</span>
      </li>
    ),
    p: ({ children }: MarkdownComponentProps) => (
      <p className="text-gray-600 leading-relaxed mb-4">{children}</p>
    ),
    strong: ({ children }: MarkdownComponentProps) => (
      <span className="font-semibold text-gray-900 bg-blue-50/50 px-1 rounded">{children}</span>
    ),
    blockquote: ({ children }: MarkdownComponentProps) => (
      <div className="bg-amber-50 border-l-4 border-amber-400 p-4 my-4 rounded-r-lg">
        <div className="flex items-center gap-2 text-amber-800 font-semibold mb-1">
          <AlertTriangle className="w-4 h-4" />
          <span>注意事项</span>
        </div>
        <div className="text-amber-900/80 italic text-sm pl-6">{children}</div>
      </div>
    ),
  };

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      {/* Header Section */}
      <div className="bg-gradient-to-r from-blue-50 to-indigo-50 px-6 py-5 border-b border-blue-100/50">
        <div className="flex items-center gap-3 mb-2">
          <div className="p-2 bg-white rounded-lg shadow-sm text-blue-600">
            <Bot className="w-6 h-6" />
          </div>
          <div>
            <h2 className="text-lg font-bold text-gray-900">AI 智能购买建议</h2>
            <p className="text-sm text-gray-500 flex items-center gap-1">
              <Sparkles className="w-3 h-3 text-amber-500" />
              基于 {rankingsList?.length || 0} 个品牌分析生成的深度建议
            </p>
          </div>
        </div>
      </div>

      {/* Content Section */}
      <div className="p-6">
        <div className="prose prose-blue max-w-none">
          <ReactMarkdown 
            components={markdownComponents}
          >
            {recommendation}
          </ReactMarkdown>
        </div>

        {/* Footer Tags - Visual decoration based on content categories */}
        <div className="mt-8 pt-6 border-t border-gray-100 flex flex-wrap gap-3">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-green-50 text-green-700 text-xs font-medium border border-green-100">
            <CheckCircle2 className="w-3.5 h-3.5" />
            值得买
          </div>
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-blue-50 text-blue-700 text-xs font-medium border border-blue-100">
            <Target className="w-3.5 h-3.5" />
            分场景推荐
          </div>
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-amber-50 text-amber-700 text-xs font-medium border border-amber-100">
            <AlertTriangle className="w-3.5 h-3.5" />
            避坑指南
          </div>
        </div>
      </div>
    </div>
  );
};

export default EnhancedSummary;
