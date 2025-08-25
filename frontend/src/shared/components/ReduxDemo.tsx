import React from 'react';
import { useAppSelector, useAppDispatch } from '../hooks/redux';
import { showSuccess, showError, showInfo } from '../store/uiSlice';
import { useGetMessagesQuery, useCreateMessageMutation } from '../../features/dashboard/store/dashboardApi';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Input } from './ui/input';
import { Loader2, MessageSquare, Plus } from 'lucide-react';

// This component demonstrates the new Redux Toolkit + RTK Query patterns
export const ReduxDemo: React.FC = () => {
  const dispatch = useAppDispatch();
  
  // RTK Query hooks for API calls
  const { data: messages, isLoading, error, refetch } = useGetMessagesQuery();
  const [createMessage, { isLoading: isCreating }] = useCreateMessageMutation();
  
  // Local state for new message
  const [newMessage, setNewMessage] = React.useState('');
  
  // Redux state selectors
  const theme = useAppSelector((state) => state.theme.theme);
  const notifications = useAppSelector((state) => state.ui.notifications);
  
  const handleCreateMessage = async () => {
    if (!newMessage.trim()) return;
    
    try {
      await createMessage({ content: newMessage }).unwrap();
      setNewMessage('');
      dispatch(showSuccess('Message created successfully!'));
    } catch (err: any) {
      dispatch(showError('Failed to create message: ' + (err.data?.message || err.message)));
    }
  };
  
  const handleShowNotification = (type: 'success' | 'error' | 'info') => {
    switch (type) {
      case 'success':
        dispatch(showSuccess('This is a success notification!'));
        break;
      case 'error':
        dispatch(showError('This is an error notification!'));
        break;
      case 'info':
        dispatch(showInfo('This is an info notification!'));
        break;
    }
  };

  return (
    <Card className="dark:bg-gray-800 dark:border-gray-700">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 dark:text-white">
          <MessageSquare className="w-5 h-5" />
          Redux Toolkit + RTK Query Demo
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* Theme Display */}
        <div className="p-3 bg-gray-100 dark:bg-gray-700 rounded-lg">
          <p className="text-sm text-gray-600 dark:text-gray-300">
            Current theme: <span className="font-semibold">{theme}</span>
          </p>
          <p className="text-xs text-gray-500 dark:text-gray-400">
            Notifications: {notifications.length}
          </p>
        </div>

        {/* Message Creation */}
        <div className="flex gap-2">
          <Input
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Enter a message..."
            className="flex-1"
          />
          <Button 
            onClick={handleCreateMessage} 
            disabled={isCreating || !newMessage.trim()}
            size="sm"
          >
            {isCreating ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <Plus className="w-4 h-4" />
            )}
          </Button>
        </div>

        {/* Notification Buttons */}
        <div className="flex gap-2">
          <Button 
            onClick={() => handleShowNotification('success')} 
            variant="outline" 
            size="sm"
            className="text-green-600 border-green-200 hover:bg-green-50"
          >
            Show Success
          </Button>
          <Button 
            onClick={() => handleShowNotification('error')} 
            variant="outline" 
            size="sm"
            className="text-red-600 border-red-200 hover:bg-red-50"
          >
            Show Error
          </Button>
          <Button 
            onClick={() => handleShowNotification('info')} 
            variant="outline" 
            size="sm"
            className="text-blue-600 border-blue-200 hover:bg-blue-50"
          >
            Show Info
          </Button>
        </div>

        {/* Messages Display */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <h4 className="font-medium dark:text-white">Messages ({messages?.length || 0})</h4>
            <Button onClick={() => refetch()} variant="outline" size="sm">
              Refresh
            </Button>
          </div>
          
          {isLoading ? (
            <div className="flex items-center justify-center p-4">
              <Loader2 className="w-6 h-6 animate-spin" />
            </div>
          ) : error ? (
            <div className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
              <p className="text-sm text-red-600 dark:text-red-400">
                Error loading messages: {(error as any)?.data?.message || 'Unknown error'}
              </p>
            </div>
          ) : messages && messages.length > 0 ? (
            <div className="space-y-2 max-h-40 overflow-y-auto">
              {messages.map((message: any) => (
                <div key={message.id} className="p-2 bg-gray-50 dark:bg-gray-700 rounded text-sm">
                  <p className="dark:text-gray-300">{message.content}</p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {new Date(message.created_at).toLocaleString()}
                  </p>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-gray-500 dark:text-gray-400 text-center p-4">
              No messages yet. Create one above!
            </p>
          )}
        </div>

        {/* Redux Benefits */}
        <div className="p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg">
          <h5 className="font-medium text-blue-800 dark:text-blue-200 mb-2">Redux Benefits:</h5>
          <ul className="text-xs text-blue-700 dark:text-blue-300 space-y-1">
            <li>• Automatic caching and background updates</li>
            <li>• Loading and error states handled automatically</li>
            <li>• Optimistic updates and rollbacks</li>
            <li>• Centralized state management</li>
            <li>• TypeScript support throughout</li>
          </ul>
        </div>
      </CardContent>
    </Card>
  );
}; 