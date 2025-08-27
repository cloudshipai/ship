"""
Sample Python code for OpenCode testing
"""
import os
import json
from typing import List, Dict, Optional

class DataProcessor:
    """Processes data from various sources"""
    
    def __init__(self, config_path: str):
        self.config_path = config_path
        self.config = self._load_config()
        
    def _load_config(self) -> Dict:
        """Load configuration from file"""
        try:
            with open(self.config_path, 'r') as f:
                return json.load(f)
        except FileNotFoundError:
            return {}
    
    def process_data(self, data: List[Dict]) -> List[Dict]:
        """Process a list of data items"""
        processed = []
        for item in data:
            if self._validate_item(item):
                processed_item = self._transform_item(item)
                processed.append(processed_item)
        return processed
    
    def _validate_item(self, item: Dict) -> bool:
        """Validate a data item"""
        required_fields = self.config.get('required_fields', [])
        return all(field in item for field in required_fields)
    
    def _transform_item(self, item: Dict) -> Dict:
        """Transform a data item"""
        # Add timestamp if not present
        if 'timestamp' not in item:
            import time
            item['timestamp'] = int(time.time())
        
        # Convert string numbers to integers
        for key, value in item.items():
            if isinstance(value, str) and value.isdigit():
                item[key] = int(value)
                
        return item

# Example usage
if __name__ == "__main__":
    processor = DataProcessor("config.json")
    
    sample_data = [
        {"id": "1", "name": "Alice", "score": "95"},
        {"id": "2", "name": "Bob", "score": "87"},
        {"id": "3", "name": "Charlie"}  # Missing score
    ]
    
    result = processor.process_data(sample_data)
    print(f"Processed {len(result)} items")
    for item in result:
        print(f"  {item}")