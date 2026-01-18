import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:responsible/models/foods.dart';
import 'package:flutter/foundation.dart';
import 'package:responsible/shoppinglist_page.dart';

class FoodPage extends StatefulWidget {
  const FoodPage({super.key});

  @override
  State<FoodPage> createState() => _FoodTypePageState();
}

class _FoodTypePageState extends State<FoodPage> {
  Future<List<Foods>> fetchFoods() async {
    final String url = kIsWeb 
      ? '/items' 
      : 'http://192.168.10.61:80/items';

    final response = await http.get(Uri.parse(url));

    if (response.statusCode == 200) {
      List jsonResponse = json.decode(response.body);
      return jsonResponse.map((data) => Foods.fromJson(data)).toList();
    } else {
      throw Exception('Could not load food types');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(    
      appBar: AppBar(
        title: const Text('Inventory'),
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
                Navigator.push(
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
                Navigator.pushReplacement(
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
                subtitle: Text('ID: ${item.id}'),
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