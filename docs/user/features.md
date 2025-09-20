# CMDB Lite Features Overview

CMDB Lite provides a comprehensive set of features for managing configuration items (CIs) and their relationships. This document provides a detailed overview of all features available in the application.

## Table of Contents

- [Configuration Item Management](#configuration-item-management)
- [Relationship Management](#relationship-management)
- [Graph Visualization](#graph-visualization)
- [Search and Filtering](#search-and-filtering)
- [Audit Logging](#audit-logging)
- [User Authentication and Authorization](#user-authentication-and-authorization)
- [User Interface](#user-interface)
- [Data Management](#data-management)

## Configuration Item Management

### CI Types

CMDB Lite supports several predefined CI types, and allows for custom types to be defined:

- **Servers**: Physical or virtual servers
- **Applications**: Software applications and services
- **Databases**: Database systems and instances
- **Network Devices**: Routers, switches, firewalls, etc.
- **Storage Systems**: Storage arrays and systems
- **Custom Types**: User-defined CI types to meet specific needs

### CI Attributes

Each CI can have custom attributes defined as key-value pairs:

- **System Attributes**: Automatically managed by the system (ID, creation date, last modified)
- **User-Defined Attributes**: Custom attributes specific to your organization's needs
- **Data Types**: Support for various data types including text, numbers, dates, and JSON

### CI Lifecycle Management

- **Create**: Add new CIs with detailed information
- **Read**: View CI details with all attributes and relationships
- **Update**: Modify CI information as needed
- **Delete**: Remove CIs when no longer needed
- **Archive**: Mark CIs as inactive without deleting them (future feature)

## Relationship Management

### Relationship Types

CMDB Lite supports various relationship types to model dependencies between CIs:

- **Depends On**: Indicates that one CI requires another to function
- **Hosts**: Indicates that one CI runs on another (e.g., application on a server)
- **Connects To**: Indicates network connectivity between CIs
- **Part Of**: Indicates that one CI is a component of another
- **Custom Types**: User-defined relationship types for specific needs

### Relationship Management Features

- **Create Relationships**: Establish connections between CIs
- **View Relationships**: See all incoming and outgoing relationships for a CI
- **Edit Relationships**: Modify relationship types or target CIs
- **Delete Relationships**: Remove relationships when no longer valid
- **Impact Analysis**: Understand the potential impact of changes based on relationships

## Graph Visualization

### Interactive Graph

The graph visualization provides an interactive way to explore relationships between CIs:

- **Force-Directed Layout**: Automatically arranges CIs to minimize overlap
- **Zoom and Pan**: Navigate through large graphs with ease
- **Node Selection**: Click on nodes to see detailed information
- **Expand/Collapse**: Focus on specific areas of the graph by expanding or collapsing nodes
- **Color Coding**: Different colors for different CI types for easy identification

### Graph Features

- **Filter by CI Type**: Show only specific types of CIs in the graph
- **Search and Highlight**: Find specific CIs and highlight them in the graph
- **Layout Options**: Choose from different graph layout algorithms
- **Export**: Export graphs as images for documentation (future feature)
- **Full-Screen Mode**: View the graph in full-screen for better visibility

## Search and Filtering

### Search Functionality

- **Global Search**: Search across all CIs from any page
- **Attribute Search**: Search within specific CI attributes
- **Fuzzy Search**: Find CIs even with approximate spelling
- **Search History**: Keep track of recent searches (future feature)

### Filtering Options

- **Filter by CI Type**: Show only specific types of CIs
- **Filter by Tags**: Show CIs with specific tags
- **Filter by Attributes**: Filter CIs based on attribute values
- **Date Range Filtering**: Filter CIs based on creation or modification dates
- **Saved Filters**: Save frequently used filter combinations (future feature)

### Sorting Options

- **Sort by Name**: Alphabetical sorting of CIs
- **Sort by Type**: Group CIs by type
- **Sort by Date**: Sort by creation or modification date
- **Custom Sorting**: Define custom sorting criteria (future feature)

## Audit Logging

### Comprehensive Logging

All actions in CMDB Lite are logged for audit and troubleshooting purposes:

- **CI Actions**: Log all create, update, and delete actions on CIs
- **Relationship Actions**: Log all create, update, and delete actions on relationships
- **User Actions**: Log login, logout, and password change actions
- **System Actions**: Log system configuration changes

### Log Features

- **Detailed Information**: Each log entry includes who, what, when, and how
- **Change Tracking**: Track before and after values for updates
- **Filtering**: Filter logs by user, action type, date range, and more
- **Export**: Export logs for external analysis or compliance reporting
- **Immutable Records**: Logs cannot be modified or deleted by users

### Log Retention

- **Configurable Retention**: Set how long logs are kept (system administrator setting)
- **Automatic Cleanup**: Automatically remove old logs based on retention policy
- **Archive Option**: Archive important logs before deletion (future feature)

## User Authentication and Authorization

### Authentication

- **Username/Password Login**: Traditional login with credentials
- **JWT Tokens**: Secure token-based authentication for API access
- **Session Management**: Control active sessions and logout remotely
- **Password Policy**: Enforce strong password requirements (configurable)
- **Password Reset**: Self-service password reset (future feature)

### Authorization

- **Role-Based Access Control (RBAC)**: Define roles with specific permissions
- **Predefined Roles**: 
  - **Admin**: Full access to all features and settings
  - **Editor**: Can create, update, and delete CIs and relationships
  - **Viewer**: Read-only access to CIs and relationships
- **Custom Roles**: Create custom roles with specific permissions (future feature)
- **Fine-Grained Permissions**: Control access to specific features and data

## User Interface

### Responsive Design

- **Desktop Optimized**: Full feature set on desktop browsers
- **Tablet Compatible**: Optimized layout for tablet devices
- **Mobile Friendly**: Access key features on mobile devices
- **Adaptive Layout**: Interface adapts to different screen sizes

### User Experience Features

- **Intuitive Navigation**: Clear menu structure and navigation paths
- **Contextual Help**: Help text and tooltips where needed
- **Breadcrumb Navigation**: Easy navigation back to previous pages
- **Keyboard Shortcuts**: Common actions accessible via keyboard shortcuts (future feature)
- **Dark Mode**: Alternative color scheme for reduced eye strain (future feature)

### Accessibility

- **Screen Reader Support**: Compatible with screen reader software
- **Keyboard Navigation**: Full functionality accessible via keyboard
- **High Contrast Mode**: Improved visibility for visually impaired users (future feature)
- **ARIA Labels**: Proper labeling for assistive technologies

## Data Management

### Import/Export

- **CSV Import**: Import CIs from CSV files (future feature)
- **JSON Import**: Import CIs and relationships from JSON files (future feature)
- **CSV Export**: Export CIs to CSV files for external analysis
- **JSON Export**: Export CIs and relationships to JSON for backup or migration
- **Template Export**: Export empty templates for data import (future feature)

### Data Validation

- **Required Fields**: Ensure important data is always provided
- **Data Type Validation**: Ensure data is in the correct format
- **Unique Constraints**: Prevent duplicate entries where appropriate
- **Referential Integrity**: Ensure relationships point to valid CIs

### Data Quality

- **Data Quality Reports**: Identify incomplete or inconsistent data (future feature)
- **Duplicate Detection**: Identify potentially duplicate CIs (future feature)
- **Orphaned Relationships**: Identify relationships pointing to deleted CIs
- **Data Profiling**: Analyze data completeness and accuracy (future feature)

## Integration Features

### API Access

- **RESTful API**: Full API access to all CMDB Lite features
- **Authentication**: API key or JWT token-based authentication
- **Rate Limiting**: Prevent abuse with configurable rate limits
- **API Documentation**: Interactive API documentation (future feature)

### Webhooks

- **Event Notifications**: Receive notifications when specific events occur (future feature)
- **Custom Payloads**: Configure what data is included in webhook payloads (future feature)
- **Retry Logic**: Automatic retry for failed webhook deliveries (future feature)
- **Security**: Secure webhook delivery with signature verification (future feature)

## Future Enhancements

The following features are planned for future releases:

- **CMDB Discovery**: Automatic discovery of CIs in your infrastructure
- **Change Management**: Formal change management workflows
- **Service Catalog**: Define and manage service offerings
- **SLA Management**: Track service level agreements
- **Advanced Reporting**: Custom reports and dashboards
- **Multi-tenancy**: Support for multiple organizations in a single instance
- **Plugin System**: Extend functionality with custom plugins

## Feature Requests

We welcome suggestions for new features and improvements! To request a feature:

1. Check if the feature has already been requested in our [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues)
2. If not, create a new issue with the "feature request" label
3. Provide a clear description of the feature and its benefits
4. Explain how it would fit into your workflow

For more information on contributing, see the [Contribution Guide](../project/contributing.md).