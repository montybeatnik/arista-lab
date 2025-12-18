import unittest
from unittest.mock import Mock, patch
from device_crawler import get_device_info  

class TestGetDeviceInfo(unittest.TestCase):

    @patch('builtins.print')
    def test_get_device_info_success(self, mock_print):
        # Mock the net_connect object
        net_connect = Mock()
        net_connect.send_command.side_effect = [
            "device1\n",  # Output for 'show hostname'
            "Interface lo0 is up\n"
            "  Internet address is 10.0.0.1/32\n"
        ]

        # Call the function
        hostname, loopback_ip = get_device_info(net_connect)

        # Assert the results
        self.assertEqual(hostname, "device1")
        self.assertEqual(loopback_ip, "10.0.0.1")
        mock_print.assert_not_called()

    @patch('builtins.print')
    def test_get_device_info_exception(self, mock_print):
        # Mock the net_connect object to raise an exception
        net_connect = Mock()
        net_connect.send_command.side_effect = Exception("Mocked exception")

        # Call the function
        hostname, loopback_ip = get_device_info(net_connect)

        # Assert the results
        self.assertIsNone(hostname)
        self.assertIsNone(loopback_ip)
        mock_print.assert_called_once_with("Failed to retrieve device info: Mocked exception")

    @patch('builtins.print')
    def test_get_device_info_no_loopback_ip(self, mock_print):
        # Mock the net_connect object
        net_connect = Mock()
        net_connect.send_command.side_effect = [
            "device1\n",  # Output for 'show hostname'
            "Interface lo0 is down\n"  # Output for 'show ip int lo0' without IP address
        ]

        # Call the function
        hostname, loopback_ip = get_device_info(net_connect)

        # Assert the results
        self.assertEqual(hostname, "device1")
        self.assertIsNone(loopback_ip)
        mock_print.assert_not_called()

if __name__ == '__main__':
    unittest.main()