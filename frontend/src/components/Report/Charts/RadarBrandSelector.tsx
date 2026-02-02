import React, { useState, useMemo, useCallback } from 'react';
import { BrandRadarChart } from './BrandRadarChart';
import type { ReportData } from '../../../types/report';

interface RadarBrandSelectorProps {
  data: ReportData;
  className?: string;
}

/**
 * 品牌雷达图选择器组件
 * 允许用户选择 2-4 个品牌进行对比分析
 */
export const RadarBrandSelector: React.FC<RadarBrandSelectorProps> = ({ 
  data, 
  className = '' 
}) => {
  // 初始化选中品牌，默认选中前3个（如果有），或者所有品牌（如果少于3个）
  const [selectedBrands, setSelectedBrands] = useState<string[]>(() => {
    if (!data?.brands) return [];
    return data.brands.slice(0, 3);
  });

  // 处理品牌选择变更
  const handleBrandToggle = useCallback((brand: string) => {
    setSelectedBrands(prev => {
      // 如果当前已选中，且选中数量大于2，则允许取消选中
      if (prev.includes(brand)) {
        if (prev.length <= 2) {
          // 保持至少2个品牌
          return prev;
        }
        return prev.filter(b => b !== brand);
      }
      
      // 如果当前未选中，且选中数量小于4，则允许选中
      if (prev.length >= 4) {
        return prev;
      }
      return [...prev, brand];
    });
  }, []);

  // 全选（最多选前4个）
  const handleSelectAll = useCallback(() => {
    if (!data?.brands) return;
    setSelectedBrands(data.brands.slice(0, 4));
  }, [data?.brands]);

  // 重置（恢复默认前3个）
  const handleReset = useCallback(() => {
    if (!data?.brands) return;
    setSelectedBrands(data.brands.slice(0, 3));
  }, [data?.brands]);

  // 构造传递给图表的数据对象
  // 保持原有数据结构，仅修改 brands 列表为当前选中的品牌
  const chartData = useMemo<ReportData>(() => ({
    ...data,
    brands: selectedBrands
  }), [data, selectedBrands]);

  // 如果没有数据，直接返回 null
  if (!data?.brands || data.brands.length === 0) {
    return null;
  }

  return (
    <div className={`bg-white rounded-lg shadow-sm p-4 ${className}`}>
      <div className="mb-4">
        <div className="flex justify-between items-center mb-2">
          <h3 className="text-lg font-semibold text-gray-800">品牌维度对比</h3>
          <div className="text-sm text-gray-500">
            <button 
              onClick={handleSelectAll}
              className="hover:text-blue-600 mr-3 transition-colors"
            >
              选前4个
            </button>
            <button 
              onClick={handleReset}
              className="hover:text-blue-600 transition-colors"
            >
              重置
            </button>
          </div>
        </div>
        
        <p className="text-xs text-gray-500 mb-3">
          请选择 2-4 个品牌进行对比 (当前已选: {selectedBrands.length})
        </p>

        <div className="flex flex-wrap gap-2">
          {data.brands.map(brand => {
            const isSelected = selectedBrands.includes(brand);
            const isDisabled = 
              (!isSelected && selectedBrands.length >= 4) || // 已满4个，不能再选
              (isSelected && selectedBrands.length <= 2);    // 仅剩2个，不能取消

            return (
              <label 
                key={brand}
                className={`
                  inline-flex items-center px-3 py-1.5 rounded-full text-sm cursor-pointer border transition-all select-none
                  ${isSelected 
                    ? 'bg-blue-50 border-blue-200 text-blue-700 font-medium' 
                    : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-gray-100'}
                  ${isDisabled && !isSelected ? 'opacity-50 cursor-not-allowed' : ''}
                  ${isDisabled && isSelected ? 'opacity-70 cursor-not-allowed' : ''}
                `}
              >
                <input
                  type="checkbox"
                  className="hidden"
                  checked={isSelected}
                  onChange={() => !isDisabled && handleBrandToggle(brand)}
                  disabled={isDisabled}
                />
                <span className="mr-1.5 flex items-center justify-center w-4 h-4 rounded-full border border-current text-[10px]">
                  {isSelected ? '✓' : '+'}
                </span>
                {brand}
              </label>
            );
          })}
        </div>
      </div>

      <div className="border-t border-gray-100 pt-4">
        {selectedBrands.length >= 2 ? (
          <BrandRadarChart data={chartData} />
        ) : (
          <div className="h-80 flex items-center justify-center text-gray-400 bg-gray-50 rounded">
            请至少选择2个品牌
          </div>
        )}
      </div>
    </div>
  );
};
