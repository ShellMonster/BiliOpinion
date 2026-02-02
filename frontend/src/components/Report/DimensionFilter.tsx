import React from 'react';
import type { Dimension } from '../../types/report';

// 组件属性定义
interface DimensionFilterProps {
  // 所有可用维度列表
  dimensions: Dimension[];
  // 当前选中的维度名称列表
  selectedDimensions: string[];
  // 维度选择变化时的回调函数
  onChange: (selected: string[]) => void;
}

/**
 * 维度筛选器组件
 * 允许用户选择要在报告中显示的评价维度
 */
export const DimensionFilter: React.FC<DimensionFilterProps> = ({
  dimensions,
  selectedDimensions,
  onChange,
}) => {
  // 检查是否已全选
  const isAllSelected = dimensions.length > 0 && selectedDimensions.length === dimensions.length;

  // 处理全选/取消全选
  const handleSelectAll = () => {
    if (isAllSelected) {
      // 如果已全选，则清空选择（或者保留至少一个？通常全选/全不选更直观）
      // 根据需求描述"取消全选"，这里清空
      onChange([]);
    } else {
      // 否则选择所有维度
      onChange(dimensions.map((d) => d.name));
    }
  };

  // 处理单个维度的选择切换
  const handleToggle = (dimensionName: string) => {
    if (selectedDimensions.includes(dimensionName)) {
      // 如果已选中，则移除
      onChange(selectedDimensions.filter((name) => name !== dimensionName));
    } else {
      // 如果未选中，则添加
      onChange([...selectedDimensions, dimensionName]);
    }
  };

  if (!dimensions || dimensions.length === 0) {
    return null;
  }

  return (
    <div className="bg-white p-4 rounded-lg shadow-sm border border-gray-100 mb-6">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-gray-700">评价维度筛选</h3>
        <button
          onClick={handleSelectAll}
          className="text-xs text-blue-600 hover:text-blue-800 font-medium"
        >
          {isAllSelected ? '取消全选' : '全选'}
        </button>
      </div>
      
      <div className="flex flex-wrap gap-3">
        {dimensions.map((dimension) => (
          <label
            key={dimension.name}
            className={`
              relative flex items-center px-3 py-1.5 rounded-full border text-sm cursor-pointer transition-colors
              ${
                selectedDimensions.includes(dimension.name)
                  ? 'bg-blue-50 border-blue-200 text-blue-700'
                  : 'bg-gray-50 border-gray-200 text-gray-600 hover:bg-gray-100'
              }
            `}
            title={dimension.description}
          >
            <input
              type="checkbox"
              className="sr-only" // 隐藏原生 checkbox
              checked={selectedDimensions.includes(dimension.name)}
              onChange={() => handleToggle(dimension.name)}
            />
            {/* 自定义 checkbox 外观 (可选，这里使用简单的背景色切换) */}
            <span className="select-none">{dimension.name}</span>
          </label>
        ))}
      </div>
      
      {/* 选中的维度数量提示 */}
      <div className="mt-2 text-xs text-gray-400">
        已选择 {selectedDimensions.length} / {dimensions.length} 个维度
      </div>
    </div>
  );
};
