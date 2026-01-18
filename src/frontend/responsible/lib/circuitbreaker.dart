class MyCircuitBreaker {
  int _failureCount = 0;
  final int threshold;
  final Duration resetTimeout;
  DateTime? _lastFailureTime;

  MyCircuitBreaker({this.threshold = 3, this.resetTimeout = const Duration(seconds: 30)});

  bool get isOpen {
    if (_lastFailureTime != null && 
        DateTime.now().difference(_lastFailureTime!) < resetTimeout) {
      return _failureCount >= threshold;
    }
    return false; 
  }

  void recordSuccess() {
    _failureCount = 0;
    _lastFailureTime = null;
  }

  void recordFailure() {
    _failureCount++;
    _lastFailureTime = DateTime.now();
  }
}