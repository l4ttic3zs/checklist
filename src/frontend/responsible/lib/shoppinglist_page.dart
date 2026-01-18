import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:responsible/models/foods.dart';
import 'package:responsible/foods_page.dart';
import 'package:flutter/foundation.dart';
import 'package:web_socket_channel/web_socket_channel.dart';

class ShoppingListPage extends StatefulWidget {
  const ShoppingListPage({super.key});

  @override
  State<ShoppingListPage> createState() => _ShoppingListPageState();
}

class _ShoppingListPageState extends State<ShoppingListPage> {
  late WebSocketChannel channel;

  @override
  void initState() {
    super.initState();
    final String wsUrl = kIsWeb
    ? '/ws'
    : 'ws://192.168.10.66:80/ws';

    print("Connecting: $wsUrl");
    channel = WebSocketChannel.connect(Uri.parse(wsUrl));

    channel.stream.listen(
      (message) {
        print("Message received: $message");
        _showAlert(message);
      },
      onError: (error) => print("WebSocket ERROR: $error"),
      onDone: () => print("WebSocket connection closed."),
    );
  }

  void _showAlert(String message) {
    ScaffoldMessenger.of(context).hideCurrentSnackBar();

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            const Icon(Icons.info_outline, color: Colors.white),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                message,
                style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold),
              ),
            ),
          ],
        ),
        backgroundColor: Colors.deepPurple,
        duration: const Duration(seconds: 2),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        margin: EdgeInsets.only(
          bottom: MediaQuery.of(context).size.height * 0.4,
          left: 20,
          right: 20,
        ),
        action: SnackBarAction(
          label: 'OK',
          textColor: Colors.white,
          onPressed: () {
            ScaffoldMessenger.of(context).hideCurrentSnackBar();
          },
        ),
      ),
    );
  }

  @override
  void dispose() {
    channel.sink.close();
    super.dispose();
  }

  Future<List<Foods>> fetchFoods() async {
    final String wsUrl = kIsWeb 
      ? '/shoppinglist' 
      : 'http://192.168.10.65:80/shoppinglist';

    final response = await http.get(Uri.parse(wsUrl));

    if (response.statusCode == 200) {
      List jsonResponse = json.decode(response.body);
      return jsonResponse.map((data) => Foods.fromJson(data)).toList();
    } else {
      throw Exception('Could not load food types');
    }
  }

  Future<void> _purchaseItem(String name, int count) async {
    final String purchase = kIsWeb 
        ? '/shoppinglist/purchase' 
        : 'http://192.168.10.65:80/shoppinglist/purchase';

    try {
      final response = await http.post(
        Uri.parse(purchase),
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({
          "name": name,
          "count": count,
        }),
      );

      if (response.statusCode == 200 || response.statusCode == 201) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('$name megvásárolva!')),
        );
        setState(() {});
      } else {
        throw Exception('Hiba a vásárlás során: ${response.statusCode}');
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Hiba: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Shopping List'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      drawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: [
            const DrawerHeader(
              decoration: BoxDecoration(color: Colors.deepPurple),
              child: Text('Menu', style: TextStyle(color: Colors.white, fontSize: 24)),
            ),
            ListTile(
              leading: const Icon(Icons.add_task),
              title: const Text('Shoppinglist'),
              onTap: () {
                Navigator.pop(context);
                Navigator.pushReplacement(
                  context,
                  MaterialPageRoute(builder: (context) => ShoppingListPage()),
                );
              },
            ),
            ListTile(
              leading: const Icon(Icons.add_shopping_cart),
              title: const Text('Inventory'),
              onTap: () {
                Navigator.pop(context);
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => FoodPage()),
                );
              },
            ),
          ],
        ),
      ),
      body: FutureBuilder<List<Foods>>(
        future: fetchFoods(),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          } else if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
            return const Center(child: Text('No data'));
          }

          return ListView.builder(
            itemCount: snapshot.data!.length,
            itemBuilder: (context, index) {
              final item = snapshot.data![index];
              return ListTile(
                leading: const CircleAvatar(child: Icon(Icons.restaurant)),
                title: Text(item.foodtype.name),
                subtitle: Text('Count: ${item.count}'),
                trailing: IconButton(
                  icon: const Icon(Icons.check_circle_outline, color: Colors.green),
                  onPressed: () {
                    _purchaseItem(item.foodtype.name, item.count);
                  },
                ),
                onTap: () {
                   print('Selected: ${item.foodtype.name}');
                },
              );
            },
          );
        },
      ),
    );
  }
}